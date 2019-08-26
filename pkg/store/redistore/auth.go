package redistore

import (
	"context"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/secure"
	"github.com/danikarik/okpock/pkg/store"
	"github.com/gomodule/redigo/redis"
)

func checkUser(user *api.User) error {
	if user == nil {
		return store.ErrNilStruct
	}
	if user.ID == "" {
		return store.ErrZeroID
	}
	return nil
}

// IsUsernameExists ...
func (p *Pool) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, store.ErrEmptyQueryParam
	}

	c := p.Get()
	defer c.Close()

	exists, err := redis.Bool(c.Do("EXISTS", key(khUserUsername, username)))
	if err != nil {
		return false, err
	}

	return exists, nil
}

// IsEmailExists ...
func (p *Pool) IsEmailExists(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, store.ErrEmptyQueryParam
	}

	c := p.Get()
	defer c.Close()

	exists, err := redis.Bool(c.Do("EXISTS", key(khUserEmail, email)))
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (p *Pool) saveUser(ctx context.Context, user *api.User) error {
	c := p.Get()
	defer c.Close()

	_, err := c.Do("HMSET", redis.Args{}.Add(key(khUser, user.ID)).AddFlat(user)...)
	if err != nil {
		return err
	}

	_, err = c.Do("ZADD", key(kzUsers), user.CreatedAt.Unix(), user.ID)
	if err != nil {
		return err
	}

	_, err = c.Do("SET", key(khUserUsername, user.Username), user.ID)
	if err != nil {
		return err
	}

	_, err = c.Do("SET", key(khUserEmail, user.Email), user.ID)
	if err != nil {
		return err
	}

	return nil
}

// SaveNewUser ...
func (p *Pool) SaveNewUser(ctx context.Context, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	return p.saveUser(ctx, user)
}

func (p *Pool) loadUser(ctx context.Context, id string) (*api.User, error) {
	c := p.Get()
	defer c.Close()

	data, err := redis.Values(c.Do("HGETALL", key(khUser, id)))
	if err == redis.ErrNil {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	var u = &api.User{}

	err = redis.ScanStruct(data, u)
	if err != nil {
		return nil, err
	}
	u.ID = id

	return u, nil
}

func (p *Pool) findUser(ctx context.Context, k, v string) (*api.User, error) {
	c := p.Get()
	defer c.Close()

	id, err := redis.String(c.Do("GET", key(k, v)))
	if err == redis.ErrNil {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	u, err := p.loadUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// LoadUser ...
func (p *Pool) LoadUser(ctx context.Context, id string) (*api.User, error) {
	if id == "" {
		return nil, store.ErrEmptyQueryParam
	}

	return p.loadUser(ctx, id)
}

// LoadUserByUsernameOrEmail ...
func (p *Pool) LoadUserByUsernameOrEmail(ctx context.Context, input string) (*api.User, error) {
	if input == "" {
		return nil, store.ErrEmptyQueryParam
	}

	u, err := p.findUser(ctx, khUserUsername, input)
	if err == store.ErrNotFound {
		return p.findUser(ctx, khUserEmail, input)
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// LoadUserByConfirmationToken ...
func (p *Pool) LoadUserByConfirmationToken(ctx context.Context, token string) (*api.User, error) {
	if token == "" {
		return nil, store.ErrEmptyQueryParam
	}

	return p.findUser(ctx, khUserConfirmationToken, token)
}

// LoadUserByRecoveryToken ...
func (p *Pool) LoadUserByRecoveryToken(ctx context.Context, token string) (*api.User, error) {
	if token == "" {
		return nil, store.ErrEmptyQueryParam
	}

	return p.findUser(ctx, khUserRecoveryToken, token)
}

// LoadUserByEmailChangeToken ...
func (p *Pool) LoadUserByEmailChangeToken(ctx context.Context, token string) (*api.User, error) {
	if token == "" {
		return nil, store.ErrEmptyQueryParam
	}

	return p.findUser(ctx, khUserEmailChangeToken, token)
}

// Authenticate ...
func (p *Pool) Authenticate(ctx context.Context, password string, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	ok := user.CheckPassword(password)
	if !ok {
		return store.ErrWrongPassword
	}

	now := api.Now()
	user.LastSignInAt = &now
	user.UpdatedAt = now

	return p.saveUser(ctx, user)
}

// ConfirmUser ...
func (p *Pool) ConfirmUser(ctx context.Context, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	c := p.Get()
	defer c.Close()

	now := api.Now()
	user.ConfirmedAt = &now
	user.ConfirmationToken = ""
	user.UpdatedAt = now

	return p.saveUser(ctx, user)
}

// SetConfirmationToken ...
func (p *Pool) SetConfirmationToken(ctx context.Context, confirm api.Confirmation, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	c := p.Get()
	defer c.Close()

	now := api.Now()
	user.ConfirmationToken = secure.Token()
	user.UpdatedAt = now

	if confirm == api.SignUpConfirmation {
		user.ConfirmationSentAt = &now
	} else if confirm == api.InviteConfirmation {
		user.InvitedAt = &now
	}

	_, err = c.Do("SET", key(khUserConfirmationToken, user.ConfirmationToken), user.ID)
	if err != nil {
		return err
	}

	return p.saveUser(ctx, user)
}

// RecoverUser ...
func (p *Pool) RecoverUser(ctx context.Context, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	c := p.Get()
	defer c.Close()

	now := api.Now()
	user.RecoveryToken = ""
	user.UpdatedAt = now

	return p.saveUser(ctx, user)
}

// SetRecoveryToken ...
func (p *Pool) SetRecoveryToken(ctx context.Context, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	c := p.Get()
	defer c.Close()

	now := api.Now()
	user.RecoverySentAt = &now
	user.RecoveryToken = secure.Token()
	user.UpdatedAt = now

	_, err = c.Do("SET", key(khUserRecoveryToken, user.RecoveryToken), user.ID)
	if err != nil {
		return err
	}

	return p.saveUser(ctx, user)
}

// ConfirmEmailChange ...
func (p *Pool) ConfirmEmailChange(ctx context.Context, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	c := p.Get()
	defer c.Close()

	now := api.Now()
	user.Email = user.EmailChange
	user.EmailChange = ""
	user.EmailChangeToken = ""
	user.UpdatedAt = now

	return p.saveUser(ctx, user)
}

// SetEmailChangeToken ...
func (p *Pool) SetEmailChangeToken(ctx context.Context, email string, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	c := p.Get()
	defer c.Close()

	now := api.Now()
	user.EmailChangeSentAt = &now
	user.EmailChange = email
	user.EmailChangeToken = secure.Token()

	_, err = c.Do("SET", key(khUserEmailChangeToken, user.EmailChangeToken), user.ID)
	if err != nil {
		return err
	}

	return p.saveUser(ctx, user)
}

// UpdateUsername ...
func (p *Pool) UpdateUsername(ctx context.Context, username string, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	user.Username = username
	user.UpdatedAt = api.Now()

	return p.saveUser(ctx, user)
}

// UpdatePassword ...
func (p *Pool) UpdatePassword(ctx context.Context, hash string, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	user.PasswordHash = hash
	user.UpdatedAt = api.Now()

	return p.saveUser(ctx, user)
}

// UpdateUserMetaData ...
func (p *Pool) UpdateUserMetaData(ctx context.Context, data map[string]interface{}, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	user.UserMetaData = data
	user.UpdatedAt = api.Now()

	return p.saveUser(ctx, user)
}

// UpdateAppMetaData ...
func (p *Pool) UpdateAppMetaData(ctx context.Context, data map[string]interface{}, user *api.User) error {
	err := checkUser(user)
	if err != nil {
		return err
	}

	user.AppMetaData = data
	user.UpdatedAt = api.Now()

	return p.saveUser(ctx, user)
}
