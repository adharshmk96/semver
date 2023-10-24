# `semver` - Semantic Version Management CLI Tool

semver offers a streamlined command-line experience to implement semantic versioning for your projects seamlessly. Whether you utilize git tags or prefer a .version file for non-git projects, semver has got you covered.

[![Go Report Card](https://goreportcard.com/badge/github.com/adharshmk96/semver)](https://goreportcard.com/report/github.com/adharshmk96/semver)
[![Build and Test](https://github.com/adharshmk96/semver/actions/workflows/go-build-test.yml/badge.svg)](https://github.com/adharshmk96/semver/actions/workflows/go-build-test.yml)
[![Release](https://github.com/adharshmk96/semver/actions/workflows/go-release.yml/badge.svg)](https://github.com/adharshmk96/semver/actions/workflows/go-release.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/adharshmk96/semver)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/adharshmk96/semver)

## Key Features

- Standard version format: MAJOR.MINOR.PATCH (e.g., 1.0.0)
- Supports pre-release versions with MAJOR.MINOR.PATCH-PRERELEASE format (e.g., 1.0.0-alpha.1)

## 🛠 Installation

Download Binary from [releases page](https://github.com/adharshmk96/semver/releases) and move it to `/usr/local/bin` or any `PATH` directory.


or

With Golang

```bash
go install github.com/adharshmk96/semver@latest
```

If you're facing with issues, Ensure go is setup and check GOPATH is in PATH configurations.

## 🚀 Getting Started

### `semver init <version>`

For new projects, run `semver init` to initialize the project. Optionally specify the initial version number.

```bash
semver init
```

or

```bash
semver init v1.0.0
```

> NOTE: You don't need to initialize if you're managing semver already with git tags.


### `semver get`

To get the current version of the project, run `semver get`.

```bash
semver get
```

## 📖 Managing Versions

### Release Versions

`semver major | minor | patch `

Increments the corresponding release version number by 1. 

example: 

```bash
semver major    # v1.0.0 -> v2.0.0
semver minor    # v1.0.0 -> v1.1.0
semver patch    # v1.0.0 -> v1.0.1
```

To attach a pre-release version, add `--[ alpha | beta | rc ]` flag.

example: 

```bash
semver major --alpha    # v1.0.0 -> v2.0.0-alpha.1
semver major --beta     # v1.0.0 -> v2.0.0-beta.1
semver major --rc       # v1.0.0 -> v2.0.0-rc.1
```

### Pre-Release Versions

`semver alpha | beta | rc`

Increments the corresponding pre-release version number by 1.

example: 

```bash
semver alpha # v1.0.0-alpha.1 -> v1.0.0-alpha.2
semver beta  # v1.0.0-beta.1 -> v1.0.0-beta.2
semver rc    # v1.0.0-rc.1 -> v1.0.0-rc.2
```

> Note: You can't run pre-release update to current release version. (ie) `semver alpha` will fail if the current version is `v1.0.0`. You have to create pre-release like for example: `semver major --alpha`.

### Remove Pre-Release

`semver release`

Removes the pre-release tag from the current version.

example: 

```bash
semver release # v1.0.0-alpha.2 -> v1.0.0
```

## 📝 Other Commands

### `semver untag <tag>`

Removes the specified tag from the git repository.

```bash
semver untag v1.0.0
```

### `semver reset`

Resets the version to specified version or `v0.0.1` ( default ).

```bash
semver reset v1.0.0 # v1.2.3 -> v1.0.0
semver reset        # v1.2.3 -> v0.0.1
```

Note: Always ensure you understand the implications of each command, especially ones like `reset`. Misuse can lead to data loss.

## License

`semver` is licensed under [MIT License](#).

## Support and Feedback

For support, questions, or feedback, please contact me at [debugslayer@gmail.com](mailto:debugslayer@gmail.com).
