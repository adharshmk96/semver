package verman

type lang string

const (
	Go lang = "go"
)

func WriteVersionConstant(version *Semver, language lang) error {
	switch language {
	case Go:
		return writeGoCmdVersionConstant(version)
	default:
		return nil
	}
}

func writeGoCmdVersionConstant(version *Semver) error {
	versionString := version.String()
	fileContent := "package cmd\n\nconst Version = \"" + versionString + "\"\n"
	filePath := "cmd/version_constant.go"

	return writeToFile(filePath, fileContent)

}
