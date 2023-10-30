package commands_test

type MockCmdExecutor struct {
	mockRunCmd func(args ...string) (string, error)
}

func (m *MockCmdExecutor) RunCmd(args ...string) (string, error) {
	return m.mockRunCmd(args...)
}
