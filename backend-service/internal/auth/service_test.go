package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ── Mock repository ───────────────────────────────────────────────────────────

type mockRepo struct {
	fnGetUserByUsername       func(context.Context, string) (*User, error)
	fnGetUserRoles            func(context.Context, int64) ([]string, error)
	fnGetUserWithOrgById      func(context.Context, int64) (*MeResponse, error)
	fnGetPasswordHashByUserId func(context.Context, int64) (string, error)
	fnUpdateMyProfile         func(context.Context, int64, *string, *string, *string, *string) error
	fnCreateUser              func(context.Context, string, string, string, *string, *int64) (int64, error)
}

func (m *mockRepo) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	if m.fnGetUserByUsername != nil {
		return m.fnGetUserByUsername(ctx, username)
	}
	return nil, ErrUserNotFound
}

func (m *mockRepo) GetUserRoles(ctx context.Context, userId int64) ([]string, error) {
	if m.fnGetUserRoles != nil {
		return m.fnGetUserRoles(ctx, userId)
	}
	return nil, nil
}

func (m *mockRepo) GetUserWithOrgById(ctx context.Context, userId int64) (*MeResponse, error) {
	if m.fnGetUserWithOrgById != nil {
		return m.fnGetUserWithOrgById(ctx, userId)
	}
	return &MeResponse{UserId: userId, Role: "CLIENT", Status: "active"}, nil
}

func (m *mockRepo) GetPasswordHashByUserId(ctx context.Context, userId int64) (string, error) {
	if m.fnGetPasswordHashByUserId != nil {
		return m.fnGetPasswordHashByUserId(ctx, userId)
	}
	return "", nil
}

func (m *mockRepo) UpdateMyProfile(ctx context.Context, userId int64, fullName, email, phone *string, passwordHash *string) error {
	if m.fnUpdateMyProfile != nil {
		return m.fnUpdateMyProfile(ctx, userId, fullName, email, phone, passwordHash)
	}
	return nil
}

func (m *mockRepo) CreateUser(ctx context.Context, username, passwordHash, fullName string, email *string, organizationId *int64) (int64, error) {
	if m.fnCreateUser != nil {
		return m.fnCreateUser(ctx, username, passwordHash, fullName, email, organizationId)
	}
	return 1, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func makeTestToken(t *testing.T, secret string, userId int64, expiresAt time.Time) string {
	t.Helper()
	claims := &JWTClaims{
		UserId:   userId,
		Username: "testuser",
		Role:     "CLIENT",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gate-backoffice",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("makeTestToken: %v", err)
	}
	return str
}

func newService(repo IRepository) IService {
	return NewService(repo, "test-secret", 24)
}

// ── ValidateToken ─────────────────────────────────────────────────────────────

func TestValidateToken_ValidToken(t *testing.T) {
	svc := newService(&mockRepo{})
	tokenStr := makeTestToken(t, "test-secret", 42, time.Now().Add(time.Hour))

	claims, err := svc.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.UserId != 42 {
		t.Errorf("expected userId=42, got %d", claims.UserId)
	}
	if claims.Role != "CLIENT" {
		t.Errorf("expected role=CLIENT, got %s", claims.Role)
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	svc := newService(&mockRepo{})
	tokenStr := makeTestToken(t, "test-secret", 1, time.Now().Add(-time.Hour))

	_, err := svc.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	svc := newService(&mockRepo{})
	tokenStr := makeTestToken(t, "other-secret", 1, time.Now().Add(time.Hour))

	_, err := svc.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}

func TestValidateToken_MalformedToken(t *testing.T) {
	svc := newService(&mockRepo{})

	_, err := svc.ValidateToken("not.a.jwt")
	if err == nil {
		t.Fatal("expected error for malformed token, got nil")
	}
}

// ── Login ─────────────────────────────────────────────────────────────────────

func TestLogin_ValidCredentials(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.MinCost)
	repo := &mockRepo{
		fnGetUserByUsername: func(_ context.Context, _ string) (*User, error) {
			return &User{
				UserId:        1,
				Username:      "alice",
				PasswordHash:  string(hash),
				PrimaryRoleId: "CLIENT",
				Status:        "active",
			}, nil
		},
	}
	svc := newService(repo)
	resp, err := svc.Login(context.Background(), &LoginRequest{Username: "alice", Password: "testpass"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if resp.Username != "alice" {
		t.Errorf("expected username=alice, got %s", resp.Username)
	}
	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.MinCost)
	repo := &mockRepo{
		fnGetUserByUsername: func(_ context.Context, _ string) (*User, error) {
			return &User{
				UserId:       1,
				Username:     "alice",
				PasswordHash: string(hash),
				Status:       "active",
			}, nil
		},
	}
	svc := newService(repo)
	_, err := svc.Login(context.Background(), &LoginRequest{Username: "alice", Password: "wrongpass"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_DeactivatedAccount(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.MinCost)
	repo := &mockRepo{
		fnGetUserByUsername: func(_ context.Context, _ string) (*User, error) {
			return &User{
				UserId:        1,
				Username:      "bob",
				PasswordHash:  string(hash),
				PrimaryRoleId: "CLIENT",
				Status:        "inactive",
			}, nil
		},
	}
	svc := newService(repo)
	_, err := svc.Login(context.Background(), &LoginRequest{Username: "bob", Password: "pass"})
	if !errors.Is(err, ErrAccountDeactivated) {
		t.Errorf("expected ErrAccountDeactivated, got %v", err)
	}
}

func TestLogin_NoPrimaryRole(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.MinCost)
	repo := &mockRepo{
		fnGetUserByUsername: func(_ context.Context, _ string) (*User, error) {
			return &User{
				UserId:        1,
				Username:      "carol",
				PasswordHash:  string(hash),
				PrimaryRoleId: "",
				Status:        "active",
			}, nil
		},
	}
	svc := newService(repo)
	_, err := svc.Login(context.Background(), &LoginRequest{Username: "carol", Password: "pass"})
	if err == nil {
		t.Fatal("expected error for user with no primary role, got nil")
	}
}

// ── Register ──────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	repo := &mockRepo{
		fnCreateUser: func(_ context.Context, username, _, fullName string, _ *string, _ *int64) (int64, error) {
			return 99, nil
		},
	}
	svc := newService(repo)
	resp, err := svc.Register(context.Background(), &RegisterRequest{
		Username: "newuser",
		Password: "securepass",
		FullName: "New User",
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if resp.UserId != 99 {
		t.Errorf("expected userId=99, got %d", resp.UserId)
	}
	if resp.Role != "CLIENT" {
		t.Errorf("expected role=CLIENT, got %s", resp.Role)
	}
}

func TestRegister_UsernameTaken(t *testing.T) {
	repo := &mockRepo{
		fnCreateUser: func(_ context.Context, _ string, _ string, _ string, _ *string, _ *int64) (int64, error) {
			return 0, ErrUsernameTaken
		},
	}
	svc := newService(repo)
	_, err := svc.Register(context.Background(), &RegisterRequest{
		Username: "taken",
		Password: "securepass",
		FullName: "Some User",
	})
	if !errors.Is(err, ErrUsernameTaken) {
		t.Errorf("expected ErrUsernameTaken, got %v", err)
	}
}

// ── UpdateMe ─────────────────────────────────────────────────────────────────

func TestUpdateMe_NilRequest(t *testing.T) {
	svc := newService(&mockRepo{})
	_, err := svc.UpdateMe(context.Background(), 1, nil)
	if !errors.Is(err, ErrNothingToUpdate) {
		t.Errorf("expected ErrNothingToUpdate, got %v", err)
	}
}

func TestUpdateMe_NoFieldsProvided(t *testing.T) {
	svc := newService(&mockRepo{})
	_, err := svc.UpdateMe(context.Background(), 1, &UpdateMeRequest{})
	if !errors.Is(err, ErrNothingToUpdate) {
		t.Errorf("expected ErrNothingToUpdate, got %v", err)
	}
}

func TestUpdateMe_NewPasswordTooShort(t *testing.T) {
	svc := newService(&mockRepo{})
	short := "1234567"
	cur := "current"
	_, err := svc.UpdateMe(context.Background(), 1, &UpdateMeRequest{
		NewPassword:     &short,
		CurrentPassword: &cur,
	})
	if !errors.Is(err, ErrPasswordTooShort) {
		t.Errorf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestUpdateMe_NewPasswordRequiresCurrent(t *testing.T) {
	svc := newService(&mockRepo{})
	newPass := "newpassword123"
	_, err := svc.UpdateMe(context.Background(), 1, &UpdateMeRequest{
		NewPassword: &newPass,
	})
	if !errors.Is(err, ErrPasswordChangeRequiresCurrent) {
		t.Errorf("expected ErrPasswordChangeRequiresCurrent, got %v", err)
	}
}

func TestUpdateMe_WrongCurrentPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("realpassword"), bcrypt.MinCost)
	repo := &mockRepo{
		fnGetPasswordHashByUserId: func(_ context.Context, _ int64) (string, error) {
			return string(hash), nil
		},
	}
	svc := newService(repo)
	newPass := "newpassword123"
	wrongCur := "wrongpassword"
	_, err := svc.UpdateMe(context.Background(), 1, &UpdateMeRequest{
		NewPassword:     &newPass,
		CurrentPassword: &wrongCur,
	})
	if !errors.Is(err, ErrWrongCurrentPassword) {
		t.Errorf("expected ErrWrongCurrentPassword, got %v", err)
	}
}
