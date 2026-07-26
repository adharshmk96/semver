package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidVersion is returned when a string is not a supported semver.
var ErrInvalidVersion = errors.New("invalid version format")

// preLabels are the supported pre-release labels, in precedence order.
var preLabels = []string{"alpha", "beta", "rc"}

// Semver is MAJOR.MINOR.PATCH with an optional pre-release, e.g. v1.2.3-rc.4.
// Label is "" for a release version, otherwise one of preLabels.
type Semver struct {
	Major, Minor, Patch int
	Label               string
	Pre                 int
}

func ParseSemver(version string) (*Semver, error) {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")

	core, pre, hasPre := strings.Cut(version, "-")

	nums := strings.Split(core, ".")
	if len(nums) != 3 {
		return nil, ErrInvalidVersion
	}

	s := &Semver{}
	for i, target := range []*int{&s.Major, &s.Minor, &s.Patch} {
		n, err := strconv.Atoi(nums[i])
		if err != nil {
			return nil, ErrInvalidVersion
		}
		*target = n
	}

	if !hasPre {
		return s, nil
	}

	label, num, ok := strings.Cut(pre, ".")
	if !ok || !isPreLabel(label) {
		return nil, ErrInvalidVersion
	}
	n, err := strconv.Atoi(num)
	if err != nil {
		return nil, ErrInvalidVersion
	}
	s.Label, s.Pre = label, n

	return s, nil
}

func isPreLabel(label string) bool {
	for _, l := range preLabels {
		if l == label {
			return true
		}
	}
	return false
}

func (s *Semver) String() string {
	v := fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
	if s.IsPreRelease() {
		v += fmt.Sprintf("-%s.%d", s.Label, s.Pre)
	}
	return v
}

func (s *Semver) IsPreRelease() bool { return s.Label != "" && s.Pre > 0 }

// Bump increments the given part: major, minor or patch resets any
// pre-release, a pre-release label increments (or starts) that pre-release,
// and "release" strips the pre-release.
func (s *Semver) Bump(part string) {
	switch part {
	case "major":
		s.Major++
		s.Minor, s.Patch = 0, 0
	case "minor":
		s.Minor++
		s.Patch = 0
	case "patch":
		s.Patch++
	case "alpha", "beta", "rc":
		if s.Label == part {
			s.Pre++
		} else {
			s.Label, s.Pre = part, 1
		}
		return
	}
	s.Label, s.Pre = "", 0
}
