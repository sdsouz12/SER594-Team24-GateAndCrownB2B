package auth

import "context"

// IRepository defines auth data access interface
type IRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserRoles(ctx context.Context, userId int64) ([]string, error)
	GetUserWithOrgById(ctx context.Context, userId int64) (*MeResponse, error)

	GetPasswordHashByUserId(ctx context.Context, userId int64) (string, error)
	UpdateMyProfile(ctx context.Context, userId int64, fullName, email, phone *string, passwordHash *string) error

	CreateUser(ctx context.Context, username, passwordHash, fullName string, email *string, organizationId *int64) (int64, error)
}
