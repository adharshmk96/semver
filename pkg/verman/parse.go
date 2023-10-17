package verman

import (
	"fmt"
	"strconv"
	"strings"
)

type Semver struct {
	Major int
	Minor int
	Patch int
	Alpha int
	Beta  int
	RC    int
}

const (
	Alpha = "alpha"
	Beta  = "beta"
	RC    = "rc"
)

func Parse(version string) (*Semver, error) {
	semver := &Semver{}
	var err error

	// Remove 'v' prefix if it exists.
	if version != "" && version[0] == 'v' {
		version = version[1:]
	}

	// Separate into version and pre-release parts.
	parts := strings.SplitN(version, "-", 2)
	verParts := strings.Split(parts[0], ".")

	// Parse major, minor, and patch versions.
	if len(verParts) != 3 {
		fmt.Printf("error parsing version: %s\n", version)
		return nil, ErrInvalidVersionFormat
	}

	semver.Major, err = strconv.Atoi(verParts[0])
	if err != nil {
		fmt.Printf("error parsing major version: %s\n", err)
		return nil, ErrInvalidVersionFormat
	}

	semver.Minor, err = strconv.Atoi(verParts[1])
	if err != nil {
		fmt.Printf("error parsing minor version: %s\n", err)
		return nil, ErrInvalidVersionFormat
	}

	semver.Patch, err = strconv.Atoi(verParts[2])
	if err != nil {
		fmt.Printf("error parsing patch version: %s\n", err)
		return nil, ErrInvalidVersionFormat
	}

	// Parse pre-release version.
	if len(parts) != 2 {
		return semver, nil
	}

	preRelease := parts[1]
	preParts := strings.Split(preRelease, ".")

	if len(preParts) != 2 {
		fmt.Printf("error parsing pre-release version: %s\n", preRelease)
		return nil, ErrInvalidVersionFormat
	}

	label, versionStr := preParts[0], preParts[1]
	preReleaseVersion, err := strconv.Atoi(versionStr)

	if err != nil {
		fmt.Printf("error parsing pre-release version: %s\n", err)
		return nil, ErrInvalidVersionFormat
	}

	switch label {
	case Alpha:
		semver.Alpha = preReleaseVersion
	case Beta:
		semver.Beta = preReleaseVersion
	case RC:
		semver.RC = preReleaseVersion
	default:
		fmt.Printf("unknown pre-release label: %s\n", label)
		return nil, ErrInvalidVersionFormat
	}

	return semver, nil
}

func (s *Semver) UpdateSemver(versionType string) {
	switch versionType {
	case "major":
		s.IncrementMajor()
	case "minor":
		s.IncrementMinor()
	case "patch":
		s.IncrementPatch()
	case "alpha":
		s.IncrementAlpha()
	case "beta":
		s.IncrementBeta()
	case "rc":
		s.IncrementRC()
	}
}

func (s *Semver) String() string {
	version := fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
	if s.RC > 0 {
		version += fmt.Sprintf("-%s.%d", RC, s.RC)
	}
	if s.Beta > 0 {
		version += fmt.Sprintf("-%s.%d", Beta, s.Beta)
	}
	if s.Alpha > 0 {
		version += fmt.Sprintf("-%s.%d", Alpha, s.Alpha)
	}
	return version
}

func (s *Semver) IncrementMajor() {
	s.Major++
	s.Minor = 0
	s.Patch = 0
	s.Alpha = 0
	s.Beta = 0
	s.RC = 0
}

func (s *Semver) IncrementMinor() {
	s.Minor++
	s.Patch = 0
	s.Alpha = 0
	s.Beta = 0
	s.RC = 0
}

func (s *Semver) IncrementPatch() {
	s.Patch++
	s.Alpha = 0
	s.Beta = 0
	s.RC = 0
}

func (s *Semver) IncrementAlpha() {
	s.Alpha++
	s.Beta = 0
	s.RC = 0
}

func (s *Semver) IncrementBeta() {
	s.Alpha = 0
	s.Beta++
	s.RC = 0
}

func (s *Semver) IncrementRC() {
	s.Alpha = 0
	s.Beta = 0
	s.RC++
}
