package sequel

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/secure"
	"github.com/danikarik/okpock/pkg/store"
)

func checkUser(u *api.User, opts byte) error {
	if (opts & checkNilStruct) != 0 {
		if u == nil {
			return store.ErrNilStruct
		}
	}

	if (opts & checkZeroID) != 0 {
		if u.ID == 0 {
			return store.ErrZeroID
		}
	}

	err := u.IsValid()
	if err != nil {
		return err
	}

	return nil
}

// IsUsernameExists ...
func (m *MySQL) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	query := m.builder.Select("count(1)").
		From("users").
		Where(sq.Eq{"username": username})

	cnt, err := m.countQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// IsEmailExists ...
func (m *MySQL) IsEmailExists(ctx context.Context, email string) (bool, error) {
	query := m.builder.Select("count(1)").
		From("users").
		Where(sq.Eq{"email": email})

	cnt, err := m.countQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// SaveNewUser ...
func (m *MySQL) SaveNewUser(ctx context.Context, user *api.User) error {
	err := checkUser(user, checkNilStruct)
	if err != nil {
		return err
	}

	user.Role = api.ClientRole
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := m.builder.Insert("users").
		Columns(
			"role",
			"username",
			"email",
			"password_hash",
			"raw_app_metadata",
			"raw_user_metadata",
			"created_at",
			"updated_at",
		).
		Values(
			user.Role,
			user.Username,
			user.Email,
			user.PasswordHash,
			user.AppMetaData,
			user.UserMetaData,
			user.CreatedAt,
			user.UpdatedAt,
		)

	id, err := m.insertQuery(ctx, query)
	if err != nil {
		return err
	}
	user.ID = id

	return nil
}

func (m *MySQL) loadUser(ctx context.Context, query sq.SelectBuilder) (*api.User, error) {
	row, err := m.selectRowQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	var u = &api.User{}

	err = row.StructScan(u)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// LoadUser ...
func (m *MySQL) LoadUser(ctx context.Context, id int64) (*api.User, error) {
	if id == 0 {
		return nil, store.ErrZeroID
	}

	query := m.builder.Select("*").
		From("users").
		Where(sq.Eq{"id": id})

	return m.loadUser(ctx, query)
}

// LoadUserByUsernameOrEmail ...
func (m *MySQL) LoadUserByUsernameOrEmail(ctx context.Context, input string) (*api.User, error) {
	if input == "" {
		return nil, store.ErrEmptyQueryParam
	}

	query := m.builder.Select("*").
		From("users").
		Where(sq.Or{
			sq.Eq{"username": input},
			sq.Eq{"email": input},
		})

	return m.loadUser(ctx, query)
}

// LoadUserByConfirmationToken ...
func (m *MySQL) LoadUserByConfirmationToken(ctx context.Context, token string) (*api.User, error) {
	if token == "" {
		return nil, store.ErrEmptyQueryParam
	}

	query := m.builder.Select("*").
		From("users").
		Where(sq.Eq{"confirmation_token": token})

	return m.loadUser(ctx, query)
}

// LoadUserByRecoveryToken ...
func (m *MySQL) LoadUserByRecoveryToken(ctx context.Context, token string) (*api.User, error) {
	if token == "" {
		return nil, store.ErrEmptyQueryParam
	}

	query := m.builder.Select("*").
		From("users").
		Where(sq.Eq{"recovery_token": token})

	return m.loadUser(ctx, query)
}

// LoadUserByEmailChangeToken ...
func (m *MySQL) LoadUserByEmailChangeToken(ctx context.Context, token string) (*api.User, error) {
	if token == "" {
		return nil, store.ErrEmptyQueryParam
	}

	query := m.builder.Select("*").
		From("users").
		Where(sq.Eq{"email_change_token": token})

	return m.loadUser(ctx, query)
}

// Authenticate ...
func (m *MySQL) Authenticate(ctx context.Context, password string, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	ok := user.CheckPassword(password)
	if !ok {
		return store.ErrWrongPassword
	}

	now := time.Now()
	user.LastSignInAt = &now

	query := m.builder.Update("users").
		Set("last_signin_at", user.LastSignInAt).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// ConfirmUser ...
func (m *MySQL) ConfirmUser(ctx context.Context, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.ConfirmedAt = &now
	user.ConfirmationToken = ""

	query := m.builder.Update("users").
		Set("confirmation_token", user.ConfirmationToken).
		Set("confirmed_at", user.ConfirmedAt).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetConfirmationToken ...
func (m *MySQL) SetConfirmationToken(ctx context.Context, confirm api.Confirmation, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.ConfirmationToken = secure.Token()

	query := m.builder.Update("users").
		Set("confirmation_token", user.ConfirmationToken).
		Where(sq.Eq{"id": user.ID})

	if confirm == api.SignUpConfirmation {
		user.ConfirmationSentAt = &now
		query = query.Set("confirmation_sent_at", user.ConfirmationSentAt)
	} else if confirm == api.InviteConfirmation {
		user.InvitedAt = &now
		query = query.Set("invited_at", user.InvitedAt)
	}

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// RecoverUser ...
func (m *MySQL) RecoverUser(ctx context.Context, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	user.RecoveryToken = ""

	query := m.builder.Update("users").
		Set("recovery_token", user.RecoveryToken).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetRecoveryToken ...
func (m *MySQL) SetRecoveryToken(ctx context.Context, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.RecoverySentAt = &now
	user.RecoveryToken = secure.Token()

	query := m.builder.Update("users").
		Set("recovery_token", user.RecoveryToken).
		Set("recovery_sent_at", user.RecoverySentAt).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// ConfirmEmailChange ...
func (m *MySQL) ConfirmEmailChange(ctx context.Context, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	user.Email = user.EmailChange
	user.EmailChange = ""
	user.EmailChangeToken = ""

	query := m.builder.Update("users").
		Set("email", user.Email).
		Set("email_change", user.EmailChange).
		Set("email_change_token", user.EmailChangeToken).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetEmailChangeToken ...
func (m *MySQL) SetEmailChangeToken(ctx context.Context, email string, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.EmailChangeSentAt = &now
	user.EmailChange = email
	user.EmailChangeToken = secure.Token()

	query := m.builder.Update("users").
		Set("email_change", user.EmailChange).
		Set("email_change_token", user.EmailChangeToken).
		Set("email_change_sent_at", user.EmailChangeSentAt).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUsername ...
func (m *MySQL) UpdateUsername(ctx context.Context, username string, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	user.Username = username
	user.UpdatedAt = time.Now()

	query := m.builder.Update("users").
		Set("username", user.Username).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword ...
func (m *MySQL) UpdatePassword(ctx context.Context, password string, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	hash, err := secure.NewPassword(password)
	if err != nil {
		return err
	}

	user.PasswordHash = hash
	user.UpdatedAt = time.Now()

	query := m.builder.Update("users").
		Set("password_hash", user.PasswordHash).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserMetaData ...
func (m *MySQL) UpdateUserMetaData(ctx context.Context, data map[string]interface{}, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	user.UserMetaData = data
	user.UpdatedAt = time.Now()

	query := m.builder.Update("users").
		Set("raw_user_metadata", user.UserMetaData).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// UpdateAppMetaData ...
func (m *MySQL) UpdateAppMetaData(ctx context.Context, data map[string]interface{}, user *api.User) error {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	user.AppMetaData = data
	user.UpdatedAt = time.Now()

	query := m.builder.Update("users").
		Set("raw_app_metadata", user.AppMetaData).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
