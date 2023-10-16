package verman

import (
	"os"

	"github.com/spf13/viper"
)

func GetVersionFromConfig() (*Semver, error) {
	semver := &Semver{}

	semver.Major = viper.GetInt("major")
	semver.Minor = viper.GetInt("minor")
	semver.Patch = viper.GetInt("patch")
	semver.Alpha = viper.GetInt("alpha")
	semver.Beta = viper.GetInt("beta")
	semver.RC = viper.GetInt("rc")

	return semver, nil

}

func WriteVersionToConfig(version *Semver) error {
	viper.Set("major", version.Major)
	viper.Set("minor", version.Minor)
	viper.Set("patch", version.Patch)
	viper.Set("alpha", version.Alpha)
	viper.Set("beta", version.Beta)
	viper.Set("rc", version.RC)

	return viper.WriteConfigAs(".version.yaml")
}

func RemoveConfig() error {
	return os.RemoveAll(".version.yaml")
}

func writeToFile(filePath string, fileContent string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fileContent)
	if err != nil {
		return err
	}

	return nil
}
