package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

//go:embed ui
var uiFiles embed.FS

// state is the payload the UI renders.
type state struct {
	Repo      string  `json:"repo"`
	Origin    string  `json:"origin"`
	Version   string  `json:"version"`
	Source    string  `json:"source"`
	IsGitRepo bool    `json:"isGitRepo"`
	Tags      []Tag   `json:"tags"`
	Changes   Changes `json:"changes"`
	Head      Head    `json:"head"`
}

// request is the body of every mutating API call.
type request struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Label   string `json:"label"`
	Remote  bool   `json:"remote"`
	Push    bool   `json:"push"`
	Message string `json:"message"`
	All     bool   `json:"all"`
}

// Serve starts the UI server on addr and blocks. When open is set the default
// browser is pointed at the server once it is listening.
func Serve(addr string, open bool) error {
	assets, err := fs.Sub(uiFiles, "ui")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(assets)))
	mux.HandleFunc("/api/state", handleState)
	mux.HandleFunc("/api/init", action(doInit))
	mux.HandleFunc("/api/bump", action(doBump))
	mux.HandleFunc("/api/release", action(doRelease))
	mux.HandleFunc("/api/push", action(doPush))
	mux.HandleFunc("/api/push-all", action(doPushAll))
	mux.HandleFunc("/api/delete", action(doDelete))
	mux.HandleFunc("/api/rename", action(doRename))
	mux.HandleFunc("/api/sync", action(doSync))
	mux.HandleFunc("/api/commit", action(doCommit))
	mux.HandleFunc("/api/reset", action(doReset))
	mux.HandleFunc("/api/refs", handleRefs)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	url := "http://" + listener.Addr().String()
	fmt.Println("semver ui running at", url)
	if open {
		openBrowser(url)
	}

	return http.Serve(listener, mux)
}

func handleState(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, currentState())
}

// uiContext is the project context with the version resolved to the highest
// semver tag, so the ui always bumps from the newest version rather than from
// whichever tag happens to be reachable from HEAD.
func uiContext() *Context {
	ctx := BuildContext(false)
	if ctx.IsGitRepo {
		if version := HighestVersion(); version != nil {
			ctx.Version, ctx.Source = version, SourceGit
		}
	}
	return ctx
}

func currentState() state {
	ctx := uiContext()
	s := state{
		Repo:      RepoName(),
		Origin:    OriginURL(),
		Source:    ctx.Source.String(),
		IsGitRepo: ctx.IsGitRepo,
		Tags:      []Tag{},
	}
	if ctx.Source != SourceNone {
		s.Version = ctx.Version.String()
	}
	if ctx.IsGitRepo {
		if tags, err := ListTags(); err == nil {
			s.Tags = tags
		}
		s.Changes = WorkingTree()
		s.Head = HeadCommit()
	}
	return s
}

// action wraps a mutating handler: it decodes the request, runs fn and
// responds with the refreshed state, or with the error fn returned.
func action(fn func(request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req request
		if r.ContentLength != 0 {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
		}

		if err := fn(req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, currentState())
	}
}

func doInit(req request) error {
	ctx := uiContext()
	if ctx.Source != SourceNone {
		return errors.New("semver is already initialized")
	}
	version := strings.TrimSpace(req.Version)
	if version == "" {
		version = "v0.0.1"
	}
	if _, err := ParseSemver(version); err != nil {
		return err
	}
	return ctx.Commit(version)
}

// doBump increments a part of the current version, optionally starting a
// pre-release of it, and pushes the new tag when asked.
func doBump(req request) error {
	ctx := uiContext()
	if ctx.Source == SourceNone {
		return errors.New("semver is not initialized")
	}

	part := req.Name
	isRelease := part == "major" || part == "minor" || part == "patch"
	if !isRelease && !isPreLabel(part) {
		return fmt.Errorf("unknown version part %q", part)
	}
	if !isRelease && !ctx.Version.IsPreRelease() {
		return errors.New("current version is not a pre-release. bump major, minor or patch with a pre-release label instead")
	}

	ctx.Version.Bump(part)
	if isRelease && req.Label != "" {
		if !isPreLabel(req.Label) {
			return fmt.Errorf("unknown pre-release label %q", req.Label)
		}
		ctx.Version.Bump(req.Label)
	}

	return commit(ctx, req.Push)
}

func doRelease(req request) error {
	ctx := uiContext()
	if ctx.Source == SourceNone {
		return errors.New("semver is not initialized")
	}
	if !ctx.Version.IsPreRelease() {
		return errors.New("current version is not a pre-release")
	}
	ctx.Version.Bump("release")
	return commit(ctx, req.Push)
}

func commit(ctx *Context, push bool) error {
	if err := ctx.Commit(ctx.Version.String()); err != nil {
		return err
	}
	if push {
		return ctx.Push()
	}
	return nil
}

func doPush(req request) error {
	if req.Version == "" {
		return errors.New("no version given")
	}
	return PushTag(req.Version)
}

func doPushAll(request) error { return PushAllTags() }

func doDelete(req request) error {
	if req.Version == "" {
		return errors.New("no version given")
	}
	return Untag([]string{req.Version}, req.Remote)
}

func doRename(req request) error {
	if req.Version == "" || req.Name == "" {
		return errors.New("both the current and the new version are required")
	}
	if _, err := ParseSemver(req.Name); err != nil {
		return err
	}
	return RenameTag(req.Version, req.Name, req.Remote)
}

func doSync(request) error { return FetchTags() }

func doCommit(req request) error {
	if !BuildContext(false).IsGitRepo {
		return errors.New("not a git repository")
	}
	return CommitChanges(req.Message, req.All)
}

func doReset(req request) error {
	ctx := uiContext()
	if ctx.Source == SourceNone {
		return errors.New("semver is not initialized")
	}
	return ctx.Reset(req.Remote)
}

func handleRefs(w http.ResponseWriter, r *http.Request) {
	ctx := uiContext()
	if ctx.Source == SourceNone {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "semver is not initialized"})
		return
	}
	refs, err := ctx.References()
	if err != nil {
		refs = ""
	}
	writeJSON(w, http.StatusOK, map[string]string{"version": ctx.Version.String(), "refs": refs})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func openBrowser(url string) {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "explorer"
	default:
		cmd = "xdg-open"
	}
	exec.Command(cmd, url).Start()
}
