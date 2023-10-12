package parser_test

import (
	"testing"

	"github.com/adharshmk96/semver/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseSemverString(t *testing.T) {
	t.Run("Parses correct versions.", func(t *testing.T) {

		tc := []struct {
			name     string
			input    string
			expected parser.Semver
		}{
			{
				name:  "Test correct alpha version",
				input: "v1.2.3-alpha.1",
				expected: parser.Semver{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Alpha: 1,
					Beta:  0,
					RC:    0,
				},
			},
			{
				name:  "Test correct beta version",
				input: "v1.2.3-beta.1",
				expected: parser.Semver{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Alpha: 0,
					Beta:  1,
					RC:    0,
				},
			},
			{
				name:  "Test correct rc version",
				input: "v1.2.3-rc.1",
				expected: parser.Semver{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Alpha: 0,
					Beta:  0,
					RC:    1,
				},
			},
			{
				name:  "Test correct release version",
				input: "v1.2.3",
				expected: parser.Semver{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Alpha: 0,
					Beta:  0,
					RC:    0,
				},
			},
			{
				name:  "Test correct release version without v",
				input: "1.2.3",
				expected: parser.Semver{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Alpha: 0,
					Beta:  0,
				},
			},
		}

		for _, c := range tc {
			t.Run(c.name, func(t *testing.T) {
				actual, err := parser.Parse(c.input)
				if err != nil {
					t.Errorf("Error parsing semver string: %v", err)
				}

				assert.NoError(t, err)
				assert.NotNil(t, actual)

				assert.Equal(t, c.expected.Major, actual.Major)
				assert.Equal(t, c.expected.Minor, actual.Minor)
				assert.Equal(t, c.expected.Patch, actual.Patch)
				assert.Equal(t, c.expected.Alpha, actual.Alpha)
				assert.Equal(t, c.expected.Beta, actual.Beta)
				assert.Equal(t, c.expected.RC, actual.RC)
			})
		}
	})

	t.Run("Returns error for incorrect versions.", func(t *testing.T) {
		tc := []string{
			"whtever",
			"v1",
			"v1.2",
			"v1.2.3.4",
			"v1.2.3-alpha",
			"v1.2.3-beta",
			"v1.2.3-rc",
			"v1.2.3-alpha.1.2",
			"v1.2.3-beta.1.2",
			"v1.2.3-rc.1.2",
			"v1.2.3-alpha.1-beta.1",
			"v1.2.3-alpha.1-rc.1",
			"v1.2.3-beta.1-rc.1",
			"v1.2.3-alpha.1-beta.1-rc.1",
			"v1.2.3-alpha.1-beta.1-rc.1+build.1",
			"v1.2.3-alpha.1-beta.1-rc.1+build.1.2",
		}

		for _, c := range tc {
			t.Run(c, func(t *testing.T) {
				_, err := parser.Parse(c)
				assert.Error(t, err)
			})
		}
	})
}

func TestVersionToString(t *testing.T) {
	t.Run("Returns correct string for semver.", func(t *testing.T) {
		tc := []struct {
			name     string
			input    parser.Semver
			expected string
		}{
			{
				name: "Test correct alpha version",
				input: parser.Semver{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Alpha: 1,
					Beta:  0,
					RC:    0,
				},
				expected: "v1.2.3-alpha.1",
			},
			{
				name: "Test correct beta version",
				input: parser.Semver{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Alpha: 0,
					Beta:  1,
					RC:    0,
				},
				expected: "v1.2.3-beta.1",
			},
			{
				name: "Test correct rc version",
				input: parser.Semver{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Alpha: 0,
					Beta:  0,
					RC:    1,
				},
				expected: "v1.2.3-rc.1",
			},
			{
				name: "Test correct release version",
				input: parser.Semver{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Alpha: 0,
					Beta:  0,
					RC:    0,
				},
				expected: "v1.2.3",
			},
		}

		for _, c := range tc {
			t.Run(c.name, func(t *testing.T) {
				actual := c.input.String()
				assert.Equal(t, c.expected, actual)
			})
		}
	})
}

func TestIncrement(t *testing.T) {
	t.Run("Increments major version correctly.", func(t *testing.T) {
		version := parser.Semver{
			Major: 1,
		}

		version.IncrementMajor()

		assert.Equal(t, 2, version.Major)

	})

	t.Run("Increments minor version correctly.", func(t *testing.T) {
		version := parser.Semver{
			Minor: 1,
		}

		version.IncrementMinor()

		assert.Equal(t, 2, version.Minor)

	})

	t.Run("Increments patch version correctly.", func(t *testing.T) {
		version := parser.Semver{
			Patch: 1,
		}

		version.IncrementPatch()

		assert.Equal(t, 2, version.Patch)

	})

	t.Run("Resets minor and patch version when major version is incremented.", func(t *testing.T) {
		version := parser.Semver{
			Major: 1,
			Minor: 1,
			Patch: 1,
		}

		version.IncrementMajor()

		assert.Equal(t, 2, version.Major)
		assert.Equal(t, 0, version.Minor)
		assert.Equal(t, 0, version.Patch)
	})

	t.Run("Resets alpha, beta and rc version when patch version is incremented.", func(t *testing.T) {
		version := parser.Semver{
			Patch: 1,
			Alpha: 1,
			Beta:  1,
			RC:    1,
		}

		version.IncrementPatch()

		assert.Equal(t, 2, version.Patch)
		assert.Equal(t, 0, version.Alpha)
		assert.Equal(t, 0, version.Beta)
		assert.Equal(t, 0, version.RC)
	})

	t.Run("Resets patch, alpha, beta and rc version when minor version is incremented.", func(t *testing.T) {
		version := parser.Semver{
			Minor: 1,
			Patch: 1,
			Alpha: 1,
			Beta:  1,
			RC:    1,
		}

		version.IncrementMinor()

		assert.Equal(t, 2, version.Minor)
		assert.Equal(t, 0, version.Patch)
		assert.Equal(t, 0, version.Alpha)
		assert.Equal(t, 0, version.Beta)
		assert.Equal(t, 0, version.RC)
	})

	t.Run("Resets minor, patch, alpha, beta and rc version when major version is incremented.", func(t *testing.T) {
		version := parser.Semver{
			Major: 1,
			Minor: 1,
			Patch: 1,
			Alpha: 1,
			Beta:  1,
			RC:    1,
		}

		version.IncrementMajor()

		assert.Equal(t, 2, version.Major)
		assert.Equal(t, 0, version.Minor)
		assert.Equal(t, 0, version.Patch)
		assert.Equal(t, 0, version.Alpha)
		assert.Equal(t, 0, version.Beta)
		assert.Equal(t, 0, version.RC)
	})

	t.Run("Alpha version is incremented correctly.", func(t *testing.T) {
		version := parser.Semver{
			Alpha: 1,
		}

		version.IncrementAlpha()

		assert.Equal(t, 2, version.Alpha)
	})

	t.Run("Beta version is incremented correctly.", func(t *testing.T) {
		version := parser.Semver{
			Beta: 1,
		}

		version.IncrementBeta()

		assert.Equal(t, 2, version.Beta)
	})

	t.Run("RC version is incremented correctly.", func(t *testing.T) {
		version := parser.Semver{
			RC: 1,
		}

		version.IncrementRC()

		assert.Equal(t, 2, version.RC)
	})

	t.Run("Alpha and RC version is reset when beta version is incremented.", func(t *testing.T) {
		version := parser.Semver{
			Alpha: 1,
			Beta:  1,
			RC:    1,
		}

		version.IncrementBeta()

		assert.Equal(t, 2, version.Beta)
		assert.Equal(t, 0, version.Alpha)
		assert.Equal(t, 0, version.RC)
	})

	t.Run("Alpha and Beta version is reset when rc version is incremented.", func(t *testing.T) {
		version := parser.Semver{
			Alpha: 1,
			Beta:  1,
			RC:    1,
		}

		version.IncrementRC()

		assert.Equal(t, 2, version.RC)
		assert.Equal(t, 0, version.Alpha)
		assert.Equal(t, 0, version.Beta)
	})

	t.Run("Beta and RC version is reset when alpha version is incremented.", func(t *testing.T) {
		version := parser.Semver{
			Alpha: 1,
			Beta:  1,
			RC:    1,
		}

		version.IncrementAlpha()

		assert.Equal(t, 0, version.RC)
		assert.Equal(t, 2, version.Alpha)
		assert.Equal(t, 0, version.Beta)
	})

}
