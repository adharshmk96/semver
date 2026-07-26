const LABELS = ["alpha", "beta", "rc"];

// parse turns "v1.2.3-rc.4" into its parts, or null when it is not a semver.
function parse(version) {
  const match = /^v?(\d+)\.(\d+)\.(\d+)(?:-(alpha|beta|rc)\.(\d+))?$/.exec((version || "").trim());
  if (!match) return null;
  return {
    major: +match[1],
    minor: +match[2],
    patch: +match[3],
    label: match[4] || "",
    pre: match[5] ? +match[5] : 0,
  };
}

function format(v) {
  const core = `v${v.major}.${v.minor}.${v.patch}`;
  return v.label && v.pre > 0 ? `${core}-${v.label}.${v.pre}` : core;
}

// bump mirrors the server's version bumping so the buttons can show a preview.
function bump(v, part) {
  const next = { ...v };
  if (LABELS.includes(part)) {
    if (next.label === part) next.pre++;
    else { next.label = part; next.pre = 1; }
    return next;
  }
  if (part === "major") { next.major++; next.minor = 0; next.patch = 0; }
  if (part === "minor") { next.minor++; next.patch = 0; }
  if (part === "patch") { next.patch++; }
  next.label = ""; next.pre = 0;
  return next;
}

// browsable turns an scp-style remote (git@host:owner/repo.git) into an https
// url, so the origin can be linked. Other urls are returned as they are.
function browsable(origin) {
  const match = /^(?:ssh:\/\/)?[^@\s]+@([^:/]+)[:/](.+?)(?:\.git)?$/.exec(origin || "");
  if (match) return `https://${match[1]}/${match[2]}`;
  return origin.replace(/\.git$/, "");
}

document.addEventListener("alpine:init", () => {
  Alpine.data("semver", () => ({
    state: {
      repo: "", origin: "", version: "", source: "", isGitRepo: false, tags: [],
      changes: { staged: 0, unstaged: 0, untracked: 0, files: [] },
      head: { branch: "", commit: "", subject: "", tags: [] },
    },
    loaded: false,
    busy: false,
    error: "",
    references: null,
    initVersion: "v0.0.1",
    preLabel: "",
    remote: false,
    renaming: null,
    renameTo: "",
    confirming: null,
    commitMessage: "",
    commitAll: true,

    get changed() {
      const c = this.state.changes;
      return c.staged + c.unstaged + c.untracked > 0;
    },

    // changeSummary reads like "2 staged, 1 modified, 3 untracked".
    get changeSummary() {
      const c = this.state.changes;
      return [
        [c.staged, "staged"],
        [c.unstaged, "modified"],
        [c.untracked, "untracked"],
      ].filter(([n]) => n > 0).map(([n, label]) => `${n} ${label}`).join(", ");
    },

    get isPreRelease() {
      const v = parse(this.state.version);
      return !!v && !!v.label && v.pre > 0;
    },

    get originLink() {
      return browsable(this.state.origin);
    },

    get sourceLabel() {
      if (this.state.source === "Source: Git") return "git tags";
      if (this.state.source === "Source: File") return ".version file";
      return "not initialized";
    },

    // next is the version a pre-release bump or a release would produce.
    next(part) {
      const v = parse(this.state.version);
      return v ? format(bump(v, part)) : "";
    },

    // preview is the version a bump button would produce.
    preview(part) {
      const v = parse(this.state.version);
      if (!v) return part;
      let next = bump(v, part);
      if (this.preLabel) next = bump(next, this.preLabel);
      return `${part} → ${format(next)}`;
    },

    async load() {
      await this.call("/api/state", null);
      this.loaded = true;
    },

    post(path, body) {
      return this.call(path, body || {});
    },

    // call talks to the API; every mutating endpoint answers with the fresh
    // state, so the page always renders what the repository actually holds.
    async call(path, body) {
      this.busy = true;
      this.error = "";
      try {
        const res = await fetch(path, body === null ? undefined : {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(body),
        });
        const data = await res.json();
        if (!res.ok) {
          this.error = data.error || res.statusText;
          return;
        }
        this.state = data;
      } catch (err) {
        this.error = String(err);
      } finally {
        this.busy = false;
      }
    },

    async refs() {
      this.busy = true;
      this.error = "";
      try {
        const res = await fetch("/api/refs");
        const data = await res.json();
        if (!res.ok) this.error = data.error;
        else this.references = data.refs;
      } finally {
        this.busy = false;
      }
    },

    async commit() {
      const message = this.commitMessage.trim();
      if (!message) return;
      await this.post("/api/commit", { message, all: this.commitAll });
      if (!this.error) this.commitMessage = "";
    },

    openRename(tag) {
      this.renaming = tag;
      this.renameTo = tag.name;
    },

    async submitRename() {
      const from = this.renaming.name;
      const to = this.renameTo;
      if (!to || to === from) { this.renaming = null; return; }
      this.renaming = null;
      await this.post("/api/rename", { version: from, name: to, remote: this.remote });
    },

    confirmAction(title, body, run) {
      this.confirming = { title, body, run };
    },
  }));
});
