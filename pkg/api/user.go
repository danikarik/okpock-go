package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Confirmation is an alias for confirmation type.
type Confirmation int

const (
	_ Confirmation = iota
	// SignUpConfirmation is used when `sign up` flow is used.
	SignUpConfirmation
	// InviteConfirmation is used when `invite user` flow is used.
	InviteConfirmation
)

// Role is an alias for role representation.
type Role string

const (
	// ClientRole is app user.
	ClientRole Role = "client"
	// AdminRole is app administrator.
	AdminRole Role = "admin"
	// SuperRole is super administrator.
	SuperRole Role = "super"
)

// UserField is an alias for user field.
type UserField int

const (
	_ UserField = iota
	// ConfirmationToken refers to `user.ConfirmationToken`.
	ConfirmationToken
	// RecoveryToken refers to `user.RecoveryToken`.
	RecoveryToken
	// EmailChangeToken refers to `user.EmailChangeToken`.
	EmailChangeToken
	// EmailChange refers to `user.EmailChange`.
	EmailChange
)

// NewUser returns a new instance of user.
func NewUser(username, email, pass string, userData map[string]interface{}) (*User, error) {
	hash, err := HashPassword(pass)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		UserMetaData: userData,
		AppMetaData:  map[string]interface{}{},
	}, nil
}

// User represents user row from database.
type User struct {
	ID int64 `json:"id" db:"id"`

	Role         Role       `json:"role" db:"role"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	ConfirmedAt  *time.Time `json:"confirmedAt,omitempty" db:"confirmed_at"`
	InvitedAt    *time.Time `json:"invitedAt,omitempty" db:"invited_at"`

	ConfirmationToken  sql.NullString `json:"-" db:"confirmation_token"`
	ConfirmationSentAt *time.Time     `json:"confirmationSentAt,omitempty" db:"confirmation_sent_at"`

	RecoveryToken  sql.NullString `json:"-" db:"recovery_token"`
	RecoverySentAt *time.Time     `json:"recoverySentAt,omitempty" db:"recovery_sent_at"`

	EmailChangeToken  sql.NullString `json:"-" db:"email_change_token"`
	EmailChange       sql.NullString `json:"-" db:"email_change"`
	EmailChangeSentAt *time.Time     `json:"emailChangeSentAt,omitempty" db:"email_change_sent_at"`

	LastSignInAt *time.Time `json:"lastSignInAt,omitempty" db:"last_signin_at"`

	AppMetaData  JSONMap `json:"appMetadata" db:"raw_app_metadata"`
	UserMetaData JSONMap `json:"userMetadata" db:"raw_user_metadata"`

	IsSuperAdmin bool `json:"-" db:"is_super_admin"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// IsValid checks whether input is valid or not.
func (u *User) IsValid() error {
	if u.Username == "" {
		return errors.New("username is empty")
	}
	if u.Email == "" {
		return errors.New("email is empty")
	}
	if u.PasswordHash == "" {
		return errors.New("password is empty")
	}
	return nil
}

// String returns string representation of struct.
func (u *User) String() string {
	data, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(data)
}

// IsConfirmed returns user's confirmation status.
func (u *User) IsConfirmed() bool {
	return u.ConfirmedAt != nil
}

// HasRole returns true when the users role is set to name.
func (u *User) HasRole(role Role) bool {
	return u.Role == role
}

// SetField sets string value of user field.
func (u *User) SetField(field UserField, value string) {
	switch field {
	case ConfirmationToken:
		u.ConfirmationToken.Valid = true
		u.ConfirmationToken.String = value
		break
	case RecoveryToken:
		u.RecoveryToken.Valid = true
		u.RecoveryToken.String = value
		break
	case EmailChangeToken:
		u.EmailChangeToken.Valid = true
		u.EmailChangeToken.String = value
		break
	case EmailChange:
		u.EmailChange.Valid = true
		u.EmailChange.String = value
		break
	}
}

// GetConfirmationToken is a simple wrapper for `ConfirmationToken`.
func (u *User) GetConfirmationToken() string {
	if u.ConfirmationToken.Valid {
		return u.ConfirmationToken.String
	}
	return ""
}

// GetRecoveryToken is a simple wrapper for `RecoveryToken`.
func (u *User) GetRecoveryToken() string {
	if u.RecoveryToken.Valid {
		return u.RecoveryToken.String
	}
	return ""
}

// GetEmailChangeToken is a simple wrapper for `EmailChangeToken`.
func (u *User) GetEmailChangeToken() string {
	if u.EmailChangeToken.Valid {
		return u.EmailChangeToken.String
	}
	return ""
}

// GetEmailChange is a simple wrapper for `EmailChange`.
func (u *User) GetEmailChange() string {
	if u.EmailChange.Valid {
		return u.EmailChange.String
	}
	return ""
}

// HashPassword returns the bcrypt hash of the password.
func HashPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent.
func (u *User) CheckPassword(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pass)) == nil
}
