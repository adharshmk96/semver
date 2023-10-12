package verman

import "github.com/spf13/viper"

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
