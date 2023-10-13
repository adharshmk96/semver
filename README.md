# Semver CLI

## Overview

`semver` is a Command-Line Interface (CLI) tool designed to manage your project's semantic versioning effortlessly. Utilizing a straightforward `version.yaml` file to store the current version of the project and leveraging git tags for version management.

## Installation

With Golang

```bash
go install github.com/adharshmk96/semver@latest
```

If you're facing with issues, Ensure go is setup and check GOPATH is in PATH configurations.

## Getting Started

### Initialize `semver`

Once you've installed `semver`, navigate to your project directory and run:

```bash
semver init
```

This command initializes the `semver` configuration, creating a `version.yaml` file that stores the current version of the project.
It will also try to get the latest tag from local git repository and set it as the current version if valid.

### Display the Current Project Version

To display the current version of your project, use:

```bash
semver get
```

This command retrieves the version information from the `version.yaml` file.

### Increment the Version

To increment your project's semantic version, use:

```bash
semver up (major|minor|patch|beta|alpha|rc)
```

The `up` command will automatically increment the version by one (following semantic versioning rules) and update the `version.yaml` file accordingly. Additionally, it will tag the git repository with the new version.

### Push the Current Version

To git push the current version of your project, use:

```bash
semver push
```

This command pushes your code along with the latest version tag to the remote repository. It runs `git push origin <current-version>` to push the tag to the remote repository.

### Reset Versioning (Caution)

To reset all tags and remove the semver configuration, use:

```bash
semver reset
```

**Note**: Use this command with caution as it will remove all tags and configuration file.

### Help about Commands

To get help about any command, use:

```bash
semver help [command]
```

If `[command]` is specified, it provides help for that command. If not, it shows help for `semver`.


## Flags

### Help Flag

Use the `-h` or `--help` flag with any command to display help information.

Example:

```bash
semver -h
```
or

```bash
semver --help
```

## Examples

### Initialize and Push a Version

```bash
# Initialize semver
semver init

# Increment the version
semver up

# Push the version to the remote repository
semver push
```


## License

`semver` is licensed under [MIT License](#).

## Support and Feedback

For support, questions, or feedback, please contact us at [support@email.com](mailto:dev@adharsh.com).
