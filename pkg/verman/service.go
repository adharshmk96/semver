package verman

func WriteVersionToFile(semver *Semver, filePath string) error {
	fileContent := semver.String()
	return writeToFile(filePath, fileContent)
}
