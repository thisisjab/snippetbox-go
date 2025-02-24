package mock

import "github.com/thisisjab/snippetbox-go/internal/model"

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return model.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, model.ErrInvalidCredentials
}
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
