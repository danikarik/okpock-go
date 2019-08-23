package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/danikarik/okpock/pkg/secure"
)

// Confirmation is an alias for confirmation type.
type Confirmation string

var (
	// SignUpConfirmation is used when `register` flow is initiated.
	SignUpConfirmation = Confirmation("register")
	// InviteConfirmation is used when `invite` flow is initiated.
	InviteConfirmation = Confirmation("invite")
	// RecoveryConfirmation is used when `recovery` flow is initiated.
	RecoveryConfirmation = Confirmation("recovery")
	// EmailChangeConfirmation is used when `email_change` flow is initiated.
	EmailChangeConfirmation = Confirmation("email_change")
)

// ErrUnknownConfirmation returned when no confirmation type is match.
var ErrUnknownConfirmation = errors.New("confirmation: unknown type")

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
func NewUser(username, email, hash string, userData map[string]interface{}) *User {
	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		UserMetaData: userData,
		AppMetaData:  map[string]interface{}{},
	}
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

	ConfirmationToken  string     `json:"-" db:"confirmation_token"`
	ConfirmationSentAt *time.Time `json:"confirmationSentAt,omitempty" db:"confirmation_sent_at"`

	RecoveryToken  string     `json:"-" db:"recovery_token"`
	RecoverySentAt *time.Time `json:"recoverySentAt,omitempty" db:"recovery_sent_at"`

	EmailChangeToken  string     `json:"-" db:"email_change_token"`
	EmailChange       string     `json:"-" db:"email_change"`
	EmailChangeSentAt *time.Time `json:"emailChangeSentAt,omitempty" db:"email_change_sent_at"`

	LastSignInAt *time.Time `json:"lastSignInAt,omitempty" db:"last_signin_at"`

	AppMetaData  JSONMap `json:"-" db:"raw_app_metadata"`
	UserMetaData JSONMap `json:"userMetaData" db:"raw_user_metadata"`

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

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent.
func (u *User) CheckPassword(pass string) bool {
	return secure.CheckPassword(u.PasswordHash, pass)
}
