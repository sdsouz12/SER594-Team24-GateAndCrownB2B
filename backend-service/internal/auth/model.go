package auth

import "time"

type User struct {
	UserId         int64     `json:"userId"`
	Username       string    `json:"username"`
	Email          string    `json:"email,omitempty"`
	PasswordHash   string    `json:"-"`
	FullName       string    `json:"fullName,omitempty"`
	PrimaryRoleId  string    `json:"primaryRoleId,omitempty"`
	OrganizationId *int64    `json:"organizationId,omitempty"`
	ProviderId     *int64    `json:"providerId,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Role struct {
	RoleId      string    `json:"roleId"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken    string `json:"accessToken"`
	TokenType      string `json:"tokenType"`
	ExpiresIn      int    `json:"expiresIn"`
	Role           string `json:"role"`
	UserId         int64  `json:"userId"`
	Username       string `json:"username"`
	OrganizationId *int64 `json:"organizationId,omitempty"`
}

type OrganizationInfo struct {
	OrganizationId int64  `json:"organizationId"`
	Name           string `json:"name"`
	LogoUrl        string `json:"logoUrl,omitempty"`
}

type MeResponse struct {
	UserId       int64             `json:"userId"`
	Username     string            `json:"username"`
	FullName     string            `json:"fullName,omitempty"`
	Email        string            `json:"email,omitempty"`
	Phone        string            `json:"phone,omitempty"`
	Role         string            `json:"role"`
	Status       string            `json:"status"`
	Organization *OrganizationInfo `json:"organization,omitempty"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

// UpdateMeRequest is the body for PATCH /api/auth/me. Omitted keys are left unchanged.
// To change password, send both currentPassword and newPassword (min 8 chars).
type UpdateMeRequest struct {
	FullName *string `json:"fullName"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`

	CurrentPassword *string `json:"currentPassword"`
	NewPassword     *string `json:"newPassword"`
}
