package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrAccountDeactivated = errors.New("account has been deactivated")
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserId         int64  `json:"user_id"`
	Username       string `json:"username"`
	Role           string `json:"role"`
	ProviderId     *int64 `json:"provider_id,omitempty"`
	OrganizationId *int64 `json:"organization_id,omitempty"`
	jwt.RegisteredClaims
}

// Service implements auth business logic
type Service struct {
	repo               IRepository
	jwtSecret          string
	jwtExpirationHours int
}

// NewService creates a new auth service
func NewService(repo IRepository, jwtSecret string, jwtExpirationHours int) IService {
	return &Service{
		repo:               repo,
		jwtSecret:          jwtSecret,
		jwtExpirationHours: jwtExpirationHours,
	}
}

// Login authenticates user and returns JWT token
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Get user by username
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if account is active
	if user.Status != "active" {
		return nil, ErrAccountDeactivated
	}

	// Check if user has primary role assigned
	if user.PrimaryRoleId == "" {
		return nil, fmt.Errorf("user has no primary role assigned")
	}

	// Generate JWT token
	expirationTime := time.Now().Add(time.Duration(s.jwtExpirationHours) * time.Hour)
	claims := &JWTClaims{
		UserId:         user.UserId,
		Username:       user.Username,
		Role:           user.PrimaryRoleId,
		ProviderId:     user.ProviderId,
		OrganizationId: user.OrganizationId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gate-backoffice",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		AccessToken:    tokenString,
		TokenType:      "Bearer",
		ExpiresIn:      s.jwtExpirationHours * 3600,
		Role:           user.PrimaryRoleId,
		UserId:         user.UserId,
		Username:       user.Username,
		OrganizationId: user.OrganizationId,
	}, nil
}

// GetMe returns full user info with organization for the logged-in user.
func (s *Service) GetMe(ctx context.Context, userId int64) (*MeResponse, error) {
	return s.repo.GetUserWithOrgById(ctx, userId)
}

const minNewPasswordLen = 8

// UpdateMe patches profile fields and/or password for the logged-in user.
func (s *Service) UpdateMe(ctx context.Context, userId int64, req *UpdateMeRequest) (*MeResponse, error) {
	if req == nil {
		return nil, ErrNothingToUpdate
	}

	var fullName, email, phone *string
	if req.FullName != nil {
		v := strings.TrimSpace(*req.FullName)
		fullName = &v
	}
	if req.Email != nil {
		t := strings.TrimSpace(*req.Email)
		if t == "" {
			empty := ""
			email = &empty // repository maps empty to SQL NULL
		} else {
			email = &t
		}
	}
	if req.Phone != nil {
		v := strings.TrimSpace(*req.Phone)
		phone = &v
	}

	var newHash *string
	if req.NewPassword != nil && strings.TrimSpace(*req.NewPassword) != "" {
		newPass := strings.TrimSpace(*req.NewPassword)
		if len(newPass) < minNewPasswordLen {
			return nil, ErrPasswordTooShort
		}
		if req.CurrentPassword == nil || strings.TrimSpace(*req.CurrentPassword) == "" {
			return nil, ErrPasswordChangeRequiresCurrent
		}
		cur := strings.TrimSpace(*req.CurrentPassword)
		hash, err := s.repo.GetPasswordHashByUserId(ctx, userId)
		if err != nil {
			return nil, fmt.Errorf("load user: %w", err)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(cur)); err != nil {
			return nil, ErrWrongCurrentPassword
		}
		b, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		h := string(b)
		newHash = &h
	}

	if fullName == nil && email == nil && phone == nil && newHash == nil {
		return nil, ErrNothingToUpdate
	}

	if err := s.repo.UpdateMyProfile(ctx, userId, fullName, email, phone, newHash); err != nil {
		return nil, err
	}
	return s.repo.GetUserWithOrgById(ctx, userId)
}

// Register creates a new CLIENT user account.
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	userId, err := s.repo.CreateUser(ctx, req.Username, string(hash), req.FullName, req.Email, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return &RegisterResponse{
		UserId:   userId,
		Username: req.Username,
		FullName: req.FullName,
		Role:     "CLIENT",
	}, nil
}

// ValidateToken validates JWT token and returns claims
func (s *Service) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
