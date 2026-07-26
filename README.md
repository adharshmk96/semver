# `semver` - Semantic Version Management CLI

`semver` offers a streamlined command-line experience to implement semantic versioning for your projects seamlessly. Whether you utilize git tags or prefer a `.version` file for non-git projects, `semver` has got you covered.

Go Report Card
[Build and Test](https://github.com/adharshmk96/semver/actions/workflows/go-build-test.yml)
GitHub release (latest by date)
GitHub go.mod Go version

## Key Features

- Standard version format: `MAJOR.MINOR.PATCH` (e.g., `1.0.0`)
- Supports pre-release versions with `MAJOR.MINOR.PATCH-PRERELEASE` format (e.g., `1.0.0-alpha.1`)
- Built-in web UI (`semver ui`) for every operation, served from the same single binary



## 🛠 Installation



### Homebrew (macOS)

```bash
brew install adharshmk96/tap/semver
```



### Direct Binary Download

Linux:

```bash
VERSION="1.4.0" && \
TMP_FILE="/tmp/semver_$VERSION_linux_amd64.tar.gz" && \
curl -sLo $TMP_FILE https://github.com/adharshmk96/semver/releases/download/v$VERSION/semver_$VERSION_linux_amd64.tar.gz && \
sudo tar xz -C /usr/local/bin -f $TMP_FILE semver && \
rm $TMP_FILE
```

1. Download from the [releases page](https://github.com/adharshmk96/semver/releases).
2. Decompress and Move binary to `/usr/local/bin` or any directory in your `PATH`.



### Using Golang

```bash
go install github.com/adharshmk96/semver@latest
```

**Note**: Ensure Golang is properly set up and that `GOPATH` is configured in your PATH.

## 🚀 Getting Started



### Initialize Your Project

Start a new project with the following, optionally specifying the version:

```bash
semver init        # Default initialization
semver init v1.0.0 # With v1.0.0 version
```

**Note**: Initialization is unnecessary if you're already using git tags for versioning.

### Retrieve Current Version

```bash
semver get
```



## 📖 Version Management



### Standard Release Versions

Easily increment version numbers:

```bash
semver major    # v1.0.0 -> v2.0.0
semver minor    # v1.0.0 -> v1.1.0
semver patch    # v1.0.0 -> v1.0.1
```

For pre-releases, append the desired flag:

```bash
semver major --alpha  # v1.0.0 -> v2.0.0-alpha.1
semver minor --beta   # v1.0.0 -> v1.1.0-beta.1
semver patch --rc     # v1.0.0 -> v1.0.1-rc.1
```



### Pre-Release Versions

Manage pre-release versions effortlessly:

```bash
semver alpha  # v1.0.0-alpha.1 -> v1.0.0-alpha.2
semver beta   # v1.0.0-beta.1  -> v1.0.0-beta.2
semver rc     # v1.0.0-rc.1    -> v1.0.0-rc.2
```

Tip: use `--push` or `-p` flag to push the latest tag to remote repository along with version update command.

**Note**: Direct pre-release updates on a release version will fail. First, create a pre-release as shown above.

### Transition to a Full Release

Strip pre-release tags:

```bash
semver release  # v1.0.0-alpha.2 -> v1.0.0
```



## 🖥 Web UI

Prefer clicking to typing? `semver ui` starts a local server and opens a browser
UI with everything the CLI can do:

```bash
semver ui
```

semver ui

The header shows the repository name, its origin url, the current version and
the branch and commit a new tag would land on — including a note when that
commit already carries a tag. Each semver tag is a card marked
`current` (the highest version) and `local only` when it has not been pushed
yet, with per-tag **Push**, **Edit** (rename the tag on the same commit) and
**Delete** actions.

- **major / minor / patch** create the next version, and preview it on the button. Pick a pre-release label first to get `v1.0.0-rc.1` instead of `v1.0.0`.
- **bump alpha / beta / rc** and **release** move a pre-release forward or promote it to a release.
- **Push all** pushes every local tag, **Sync** fetches tags from the remote, **References** greps the project for the current version, and **Reset** removes every tag.
- **apply to remote** makes delete, rename and reset act on origin as well; **push after tagging** pushes each new tag right away.

When the working tree is dirty, a card appears at the top: tags always point at
a commit, so it lets you commit first without leaving the page. **All** stages
everything (untracked files included) before committing, **Staged** commits only
what is already staged.

committing from the UI

In a project that has no version yet, the UI offers to initialize one:

semver ui, uninitialized project

Flags:

```bash
semver ui --port 8080   # serve on a different port (default 7420)
semver ui --no-open     # do not open a browser
```

The server binds to `127.0.0.1` only, and the assets are embedded in the
binary, so there is nothing to install.

## 📝 Auxiliary Commands



### Removing Git Tags

remove latest tag:

```bash
semver untag
```

remove specific tag:

```bash
semver untag v1.0.0
```

add `--remote` flag to remove remote tag as well.

### Pushing Git Tags

```bash
semver push
```

Pushes latest tag to remote repository. `git push origin <current tag>` is executed under the hood.

### Resetting Versions

Easily revert to a specified version or default to `v0.0.1`:

Local:

```bash
semver reset        # v1.2.3 -> v0.0.1
semver reset v1.0.0 # v1.2.3 -> v1.0.0
```

Remote:

```bash
semver reset --remote        # v1.2.3 -> v0.0.1
semver reset v1.0.0 --remote # v1.2.3 -> v1.0.0
```

⚠️ **Caution**: Understand the implications of `reset` to avoid unintended data loss.

## License

Licensed under the [MIT License](#).

## 🤝 Support & Feedback

For any support, queries, or feedback, feel free to reach out at [debugslayer@gmail.com](mailto:debugslayer@gmail.com).