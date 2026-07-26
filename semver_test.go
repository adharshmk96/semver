package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSemver(t *testing.T) {
	valid := map[string]Semver{
		"v1.2.3":         {Major: 1, Minor: 2, Patch: 3},
		"1.2.3":          {Major: 1, Minor: 2, Patch: 3},
		"v1.2.3-alpha.1": {Major: 1, Minor: 2, Patch: 3, Label: "alpha", Pre: 1},
		"v1.2.3-beta.2":  {Major: 1, Minor: 2, Patch: 3, Label: "beta", Pre: 2},
		"v1.2.3-rc.3":    {Major: 1, Minor: 2, Patch: 3, Label: "rc", Pre: 3},
	}

	for input, expected := range valid {
		t.Run(input, func(t *testing.T) {
			actual, err := ParseSemver(input)
			assert.NoError(t, err)
			assert.Equal(t, &expected, actual)
		})
	}

	invalid := []string{
		"whatever", "v1", "v1.2", "v1.2.3.4", "v1.2.3-alpha", "v1.2.3-rc",
		"v1.2.3-alpha.1.2", "v1.2.3-alpha.1-beta.1", "v1.2.3-nightly.1",
		"v1.2.3-alpha.1+build.1",
	}

	for _, input := range invalid {
		t.Run(input, func(t *testing.T) {
			_, err := ParseSemver(input)
			assert.ErrorIs(t, err, ErrInvalidVersion)
		})
	}
}

func TestString(t *testing.T) {
	for _, input := range []string{"v1.2.3", "v1.2.3-alpha.1", "v1.2.3-beta.1", "v1.2.3-rc.1"} {
		version, err := ParseSemver(input)
		assert.NoError(t, err)
		assert.Equal(t, input, version.String())
	}
}

func TestBump(t *testing.T) {
	tc := []struct {
		from, part, expected string
	}{
		{"v1.2.3", "major", "v2.0.0"},
		{"v1.2.3", "minor", "v1.3.0"},
		{"v1.2.3", "patch", "v1.2.4"},
		{"v1.2.3-rc.1", "major", "v2.0.0"},
		{"v1.2.3-rc.1", "minor", "v1.3.0"},
		{"v1.2.3-rc.1", "patch", "v1.2.4"},
		{"v1.2.3", "alpha", "v1.2.3-alpha.1"},
		{"v1.2.3-alpha.1", "alpha", "v1.2.3-alpha.2"},
		{"v1.2.3-alpha.1", "beta", "v1.2.3-beta.1"},
		{"v1.2.3-beta.1", "rc", "v1.2.3-rc.1"},
		{"v1.2.3-rc.2", "release", "v1.2.3"},
		{"v1.2.3", "release", "v1.2.3"},
	}

	for _, c := range tc {
		t.Run(c.from+" "+c.part, func(t *testing.T) {
			version, err := ParseSemver(c.from)
			assert.NoError(t, err)

			version.Bump(c.part)
			assert.Equal(t, c.expected, version.String())
		})
	}
}

func TestIsPreRelease(t *testing.T) {
	preReleases := []string{"v0.0.0-alpha.1", "v1.2.3-beta.1", "v1.2.3-rc.1"}
	for _, input := range preReleases {
		version, _ := ParseSemver(input)
		assert.True(t, version.IsPreRelease(), input)
	}

	releases := []string{"v0.0.0", "v1.0.0", "v1.2.3"}
	for _, input := range releases {
		version, _ := ParseSemver(input)
		assert.False(t, version.IsPreRelease(), input)
	}
}
