package auth

import "errors"

var (
	// ErrWrongCurrentPassword is returned when changing password and current password does not match.
	ErrWrongCurrentPassword = errors.New("current password is incorrect")
	// ErrPasswordChangeRequiresCurrent is returned when newPassword is set but currentPassword is missing.
	ErrPasswordChangeRequiresCurrent = errors.New("currentPassword is required to set newPassword")
	// ErrPasswordTooShort is returned when newPassword is shorter than the minimum length.
	ErrPasswordTooShort = errors.New("newPassword must be at least 8 characters")
	// ErrNothingToUpdate is returned when no updatable fields were provided.
	ErrNothingToUpdate = errors.New("no fields to update")
	// ErrDuplicateEmail is returned when email is already used by another user.
	ErrDuplicateEmail = errors.New("email is already taken")
	// ErrUsernameTaken is returned when registering with an already-used username.
	ErrUsernameTaken = errors.New("username is already taken")
)
