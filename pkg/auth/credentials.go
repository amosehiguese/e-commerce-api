package auth

import "fmt"

// GetRoleCredentials maps a role to its corresponding set of credentials.
func GetRoleCredentials(role Role) ([]string, error) {
	switch role {
	case UserRole:
		return []string{
			OrderCreateCredential,
			OrderReadCredential,
			OrderCancelCredential,
		}, nil
	case AdminRole:
		return []string{
			ProductCreateCredential,
			ProductReadCredential,
			ProductUpdateCredential,
			ProductDeleteCredential,
			OrderReadCredential,
			OrderCreateCredential,
			OrderUpdateCredential,
			OrderCancelCredential,
		}, nil
	default:
		return nil, fmt.Errorf("role '%v' does not exist", role)
	}
}
