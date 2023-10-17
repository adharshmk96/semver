# `semver` - Semantic Version Management CLI Tool

`semver` is a command-line interface tool designed to simplify and streamline the management of your project's semantic versioning. It uses a `.version.yaml` file to store the project's current version and also uses git tags for version management. When updating the version, `semver` will also tag the git repository with the new version.

[![Build and Test](https://github.com/adharshmk96/semver/actions/workflows/go-build-test.yml/badge.svg)](https://github.com/adharshmk96/semver/actions/workflows/go-build-test.yml)
[![Release](https://github.com/adharshmk96/semver/actions/workflows/go-release.yml/badge.svg)](https://github.com/adharshmk96/semver/actions/workflows/go-release.yml)


## 🛠 Installation

With Golang

```bash
go install github.com/adharshmk96/semver@latest
```

If you're facing with issues, Ensure go is setup and check GOPATH is in PATH configurations.

or

Download Binary from [releases page](https://github.com/adharshmk96/semver/releases) and move it to `/usr/local/bin`

## 🚀 Usage

```
semver [command] [arguments]
```

## 📝 Commands

| Command             | Description                                                                                                                                                               |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `get`               | Retrieve and display the current version of the project.                                                                                                                  |
| `help`              | Provides detailed information and guidance on any specified command.                                                                                                      |
| `init`              | Kickstarts `semver` by initializing the necessary configuration. This includes the creation of the `.version.yaml` file. ( optionally add version )                       |
| `push`              | Pushes the current version of the project to the git repository. This helps in keeping the repository synchronized with version changes.                                  |
| `reset`             | ⚠️ (Use with caution) Wipes out all existing git tags and deletes the semver configuration. When used with `-r` or `--remote`, this also affects the remote repository.    |
| `untag`             | Removes a specified git tag. If no tag is specified, the current tag is deleted by default. Using `-r` or `--remote` will also delete the tag from the remote repository. |
| `up <version-type>` | Elevates the semver version by a single increment, adhering to semantic versioning standards. `version-type: <major/minor/patch/alpha/beta/rc>`                           |
| `version`           | Displays the running version of the `semver` CLI tool itself.                                                                                                             |

Note: Always ensure you understand the implications of each command, especially ones like `reset`. Misuse can lead to data loss.

## License

`semver` is licensed under [MIT License](#).

## Support and Feedback

For support, questions, or feedback, please contact me at [dev@adharsh.in](mailto:dev@adharsh.in).
