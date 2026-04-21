package auth

import "context"

// IService defines auth business logic interface
type IService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	GetMe(ctx context.Context, userId int64) (*MeResponse, error)
	UpdateMe(ctx context.Context, userId int64, req *UpdateMeRequest) (*MeResponse, error)
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
}
