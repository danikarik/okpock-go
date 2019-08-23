package api

import "context"

// Auth implements user related spec.
type Auth interface {
	// IsUsernameExists ...
	// TODO: description
	IsUsernameExists(ctx context.Context, username string) (bool, error)

	// IsEmailExists ...
	// TODO: description
	IsEmailExists(ctx context.Context, email string) (bool, error)

	// SaveNewUser ...
	// TODO: description
	SaveNewUser(ctx context.Context, user *User) error

	// LoadUser ...
	// TODO: description
	LoadUser(ctx context.Context, id string) (*User, error)

	// LoadUserByUsernameOrEmail ...
	// TODO: description
	LoadUserByUsernameOrEmail(ctx context.Context, input string) (*User, error)

	// LoadUserByConfirmationToken ...
	// TODO: description
	LoadUserByConfirmationToken(ctx context.Context, token string) (*User, error)

	// LoadUserByRecoveryToken ...
	// TODO: description
	LoadUserByRecoveryToken(ctx context.Context, token string) (*User, error)

	// LoadUserByEmailChangeToken ...
	// TODO: description
	LoadUserByEmailChangeToken(ctx context.Context, token string) (*User, error)

	// Authenticate ...
	// TODO: description
	Authenticate(ctx context.Context, password string, user *User) error

	// ConfirmUser ...
	// TODO: description
	ConfirmUser(ctx context.Context, user *User) error

	// SetConfirmationToken ...
	// TODO: description
	SetConfirmationToken(ctx context.Context, confirm Confirmation, user *User) error

	// RecoverUser ...
	// TODO: description
	RecoverUser(ctx context.Context, user *User) error

	// SetRecoveryToken ...
	// TODO: description
	SetRecoveryToken(ctx context.Context, user *User) error

	// ConfirmEmailChange ...
	// TODO: description
	ConfirmEmailChange(ctx context.Context, user *User) error

	// SetEmailChangeToken ...
	// TODO: description
	SetEmailChangeToken(ctx context.Context, email string, user *User) error

	// UpdateUsername ...
	// TODO: description
	UpdateUsername(ctx context.Context, username string, user *User) error

	// UpdatePassword ...
	// TODO: description
	UpdatePassword(ctx context.Context, password string, user *User) error

	// UpdateUserMetaData ...
	// TODO: description
	UpdateUserMetaData(ctx context.Context, data map[string]interface{}, user *User) error

	// UpdateAppMetaData ...
	// TODO: description
	UpdateAppMetaData(ctx context.Context, data map[string]interface{}, user *User) error
}
