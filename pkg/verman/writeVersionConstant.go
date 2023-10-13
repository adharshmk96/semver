package verman

type lang string

const (
	Go lang = "go"
)

func WriteVersionConstant(version *Semver, language lang) error {
	switch language {
	case Go:
		filePath, err := writeGoCmdVersionConstant(version)
		if err != nil {
			return err
		}

		return gitAdd(filePath)
	default:
		return nil
	}
}

func writeGoCmdVersionConstant(version *Semver) (string, error) {
	versionString := version.String()
	fileContent := "package cmd\n\nconst Version = \"" + versionString + "\"\n"
	filePath := "cmd/version_constant.go"

	return filePath, writeToFile(filePath, fileContent)

}
