# `semver` - Semantic Version Management CLI

`semver` offers a streamlined command-line experience to implement semantic versioning for your projects seamlessly. Whether you utilize git tags or prefer a `.version` file for non-git projects, `semver` has got you covered.

![Go Report Card](https://goreportcard.com/badge/github.com/adharshmk96/semver)
[![Build and Test](https://github.com/adharshmk96/semver/actions/workflows/go-build-test.yml/badge.svg)](https://github.com/adharshmk96/semver/actions/workflows/go-build-test.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/adharshmk96/semver)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/adharshmk96/semver)

## Key Features
- Standard version format: `MAJOR.MINOR.PATCH` (e.g., `1.0.0`)
- Supports pre-release versions with `MAJOR.MINOR.PATCH-PRERELEASE` format (e.g., `1.0.0-alpha.1`)

## 🛠 Installation

### Direct Binary Download

Linux:
```bash
curl -sL https://github.com/adharshmk96/semver/releases/download/v0.0.1/semver_0.0.1_linux_amd64.tar.gz | sudo tar xz -C /usr/local/bin semver
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
