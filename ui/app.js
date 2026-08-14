const LABELS = ["alpha", "beta", "rc"];
const COMMITS_PAGE = 50;

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
    tab: "version",
    loaded: false,
    busy: false,
    error: "",
    references: null,
    initVersion: "v0.0.1",
    preLabel: "",
    renaming: null,
    renameTo: "",
    renameRemote: false,
    confirming: null,
    commitMessage: "",
    commitAll: true,

    branches: { local: [], remote: [] },
    branchesLoaded: false,
    newBranchName: "",

    commits: [],
    commitsBranch: "",
    commitsSkip: 0,
    commitsHasMore: false,
    commitsLoaded: false,

    remotes: [],
    remotesLoaded: false,

    get changed() {
      const c = this.state.changes;
      return c.staged + c.unstaged + c.untracked > 0;
    },

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

    get visibleTabs() {
      const tabs = [{ id: "version", label: "Version" }];
      if (!this.state.isGitRepo) return tabs;
      tabs.push({
        id: "changes",
        label: "Changes",
        badge: this.changed ? "!" : "",
      });
      tabs.push({ id: "branches", label: "Branches" });
      tabs.push({ id: "commits", label: "Commits" });
      tabs.push({ id: "remotes", label: "Remotes" });
      return tabs;
    },

    next(part) {
      const v = parse(this.state.version);
      return v ? format(bump(v, part)) : "";
    },

    preview(part) {
      const v = parse(this.state.version);
      if (!v) return part;
      let next = bump(v, part);
      if (this.preLabel) next = bump(next, this.preLabel);
      return `${part} → ${format(next)}`;
    },

    async load() {
      await this.fetchState();
      this.loaded = true;
      await this.refreshTabData();
    },

    async localRefresh() {
      await this.fetchState();
      await this.refreshTabData();
    },

    async remoteSync() {
      await this.post("/api/sync");
      await this.refreshTabData();
    },

    async switchTab(id) {
      this.tab = id;
      await this.refreshTabData();
    },

    async refreshTabData() {
      if (!this.state.isGitRepo) return;
      if (this.tab === "branches") await this.loadBranches();
      if (this.tab === "commits") await this.loadCommits();
      if (this.tab === "remotes") await this.loadRemotes();
    },

    async fetchState() {
      this.busy = true;
      this.error = "";
      try {
        const res = await fetch("/api/state");
        const data = await res.json();
        if (!res.ok) {
          this.error = data.error || res.statusText;
          return;
        }
        this.state = data;
        if (!this.state.isGitRepo && this.tab !== "version") {
          this.tab = "version";
        }
      } catch (err) {
        this.error = String(err);
      } finally {
        this.busy = false;
      }
    },

    post(path, body) {
      return this.call(path, body || {});
    },

    // call talks to the API; mutating endpoints answer with the fresh state.
    async call(path, body) {
      this.busy = true;
      this.error = "";
      try {
        const res = await fetch(path, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(body),
        });
        const data = await res.json();
        if (!res.ok) {
          this.error = data.error || res.statusText;
          return false;
        }
        this.state = data;
        return true;
      } catch (err) {
        this.error = String(err);
        return false;
      } finally {
        this.busy = false;
      }
    },

    async fetchJSON(path) {
      this.busy = true;
      this.error = "";
      try {
        const res = await fetch(path);
        const data = await res.json();
        if (!res.ok) {
          this.error = data.error || res.statusText;
          return null;
        }
        return data;
      } catch (err) {
        this.error = String(err);
        return null;
      } finally {
        this.busy = false;
      }
    },

    async loadBranches() {
      const data = await this.fetchJSON("/api/branches");
      if (!data) return;
      this.branches = { local: data.local || [], remote: data.remote || [] };
      this.branchesLoaded = true;
    },

    async loadCommits() {
      const branch = this.state.head.branch;
      const data = await this.fetchJSON(
        `/api/commits?limit=${COMMITS_PAGE}&skip=${this.commitsSkip}&branch=${encodeURIComponent(branch)}`
      );
      if (!data) return;
      this.commits = data.commits || [];
      this.commitsBranch = data.branch || branch;
      this.commitsHasMore = this.commits.length === COMMITS_PAGE;
      this.commitsLoaded = true;
    },

    async nextCommits() {
      this.commitsSkip += COMMITS_PAGE;
      await this.loadCommits();
    },

    async prevCommits() {
      this.commitsSkip = Math.max(0, this.commitsSkip - COMMITS_PAGE);
      await this.loadCommits();
    },

    async loadRemotes() {
      const data = await this.fetchJSON("/api/remotes");
      if (!data) return;
      this.remotes = data.remotes || [];
      this.remotesLoaded = true;
    },

    async checkoutBranch(name, isRemote) {
      const ok = await this.call("/api/checkout", { branch: name, remote: isRemote });
      if (ok) {
        this.commitsSkip = 0;
        await this.refreshTabData();
      }
    },

    async deleteBranch(name, isRemote) {
      const ok = await this.call("/api/branch-delete", { branch: name, remote: isRemote });
      if (ok) await this.loadBranches();
    },

    async createBranch() {
      const name = this.newBranchName.trim();
      if (!name) return;
      const ok = await this.call("/api/branch-create", { branch: name });
      if (ok) {
        this.newBranchName = "";
        await this.loadBranches();
      }
    },

    async fetchRemote(name) {
      const ok = await this.call("/api/fetch-remote", { name });
      if (ok) {
        await this.loadRemotes();
        if (this.tab === "branches") await this.loadBranches();
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
      const ok = await this.post("/api/commit", { message, all: this.commitAll });
      if (ok) {
        this.commitMessage = "";
        await this.refreshTabData();
      }
    },

    openRename(tag) {
      this.renaming = tag;
      this.renameTo = tag.name;
      this.renameRemote = false;
    },

    async submitRename() {
      const from = this.renaming.name;
      const to = this.renameTo;
      if (!to || to === from) { this.renaming = null; return; }
      this.renaming = null;
      await this.post("/api/rename", { version: from, name: to, remote: this.renameRemote });
    },

    confirmAction(title, body, run) {
      this.confirming = { title, body, run };
    },
  }));
});
