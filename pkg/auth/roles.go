package auth

import "fmt"

type Role string

func (r Role) String() string {
	return string(r)
}

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

func VerifyRole(rl string) (Role, error) {
	switch rl {
	case string(UserRole):
		return UserRole, nil
	case string(AdminRole):
		return AdminRole, nil
	default:
		return "", fmt.Errorf("role '%v' does not exist", rl)
	}
}
