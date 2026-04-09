package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository implements auth data access
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new auth repository
func NewRepository(db *pgxpool.Pool) IRepository {
	return &Repository{db: db}
}

// GetUserByUsername retrieves user by username
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT user_id, username, email, password_hash, full_name, 
		       primary_role_id, organization_id, provider_id, status, created_at, updated_at
		FROM sys_user
		WHERE username = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.UserId,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.PrimaryRoleId,
		&user.OrganizationId,
		&user.ProviderId,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserWithOrgById retrieves full user info joined with their organization.
func (r *Repository) GetUserWithOrgById(ctx context.Context, userId int64) (*MeResponse, error) {
	var me MeResponse
	var fullName, email, phone *string
	var orgId *int64
	var orgName *string

	var orgLogo *string
	err := r.db.QueryRow(ctx, `
		SELECT u.user_id, u.username, u.full_name, u.email, u.phone_number,
		       u.primary_role_id, u.status, u.created_at, u.updated_at,
		       o.organization_id, o.name, o.logo_url
		FROM sys_user u
		LEFT JOIN bay_organization o ON o.organization_id = u.organization_id
		WHERE u.user_id = $1
	`, userId).Scan(
		&me.UserId, &me.Username, &fullName, &email, &phone,
		&me.Role, &me.Status, &me.CreatedAt, &me.UpdatedAt,
		&orgId, &orgName, &orgLogo,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if fullName != nil {
		me.FullName = *fullName
	}
	if email != nil {
		me.Email = *email
	}
	if phone != nil {
		me.Phone = *phone
	}
	if orgId != nil && orgName != nil {
		oi := &OrganizationInfo{
			OrganizationId: *orgId,
			Name:           *orgName,
		}
		if orgLogo != nil {
			oi.LogoUrl = *orgLogo
		}
		me.Organization = oi
	}
	return &me, nil
}

// GetPasswordHashByUserId returns bcrypt hash for password verification.
func (r *Repository) GetPasswordHashByUserId(ctx context.Context, userId int64) (string, error) {
	var hash string
	err := r.db.QueryRow(ctx, `SELECT password_hash FROM sys_user WHERE user_id = $1`, userId).Scan(&hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}
	return hash, nil
}

// UpdateMyProfile updates optional fields on sys_user. Empty-string email clears to NULL (unique-safe).
func (r *Repository) UpdateMyProfile(ctx context.Context, userId int64, fullName, email, phone *string, passwordHash *string) error {
	var setClauses []string
	var args []interface{}
	n := 1

	if fullName != nil {
		setClauses = append(setClauses, fmt.Sprintf("full_name = $%d", n))
		args = append(args, *fullName)
		n++
	}
	if email != nil {
		if *email == "" {
			setClauses = append(setClauses, "email = NULL")
		} else {
			setClauses = append(setClauses, fmt.Sprintf("email = $%d", n))
			args = append(args, *email)
			n++
		}
	}
	if phone != nil {
		if *phone == "" {
			setClauses = append(setClauses, "phone_number = NULL")
		} else {
			setClauses = append(setClauses, fmt.Sprintf("phone_number = $%d", n))
			args = append(args, *phone)
			n++
		}
	}
	if passwordHash != nil {
		setClauses = append(setClauses, fmt.Sprintf("password_hash = $%d", n))
		args = append(args, *passwordHash)
		n++
	}

	if len(setClauses) == 0 {
		return ErrNothingToUpdate
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, userId)
	q := fmt.Sprintf("UPDATE sys_user SET %s WHERE user_id = $%d", strings.Join(setClauses, ", "), n)

	if _, err := r.db.Exec(ctx, q, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "email") {
			return ErrDuplicateEmail
		}
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// GetUserRoles retrieves all role slugs for a user
func (r *Repository) GetUserRoles(ctx context.Context, userId int64) ([]string, error) {
	query := `
		SELECT r.slug
		FROM sys_role r
		INNER JOIN sys_user_role_rel ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating roles: %w", err)
	}

	return roles, nil
}
