package tpl

type Template struct {
	Path    string
	Content string
}

var GET_VERSION_TEMPLATE = Template{
	Path: "cmd/version_helper.go",
	Content: `package cmd

const version = "{{version}}"
`,
}
