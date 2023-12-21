package mock_script_repo

type AccessUtil struct {
	m *MockScriptAccessRepo
}

func (m *MockScriptAccessRepo) U() *AccessUtil {
	return &AccessUtil{m}
}
