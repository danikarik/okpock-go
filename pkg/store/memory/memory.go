package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/secure"
	"github.com/danikarik/okpock/pkg/store"
)

type pass struct {
	serial  string
	token   string
	id      string
	updated time.Time
}

type reg struct {
	serial string
	device string
	push   string
	id     string
}

// New returns a new instance of memory mock.
func New() *Memory {
	mock := &Memory{
		passes:           make(map[string]*pass),
		regs:             make(map[string]*reg),
		users:            make(map[int64]*api.User),
		userProjects:     make(map[int64]int64),
		projects:         make(map[int64]*api.Project),
		userUploads:      make(map[int64]int64),
		uploads:          make(map[int64]*api.Upload),
		passCards:        make(map[int64]*api.PassCardInfo),
		projectPassCards: make(map[int64]int64),
	}
	return mock
}

// Memory is mock implementor.
type Memory struct {
	mu               sync.Mutex
	passes           map[string]*pass
	regs             map[string]*reg
	users            map[int64]*api.User
	userProjects     map[int64]int64
	projects         map[int64]*api.Project
	userUploads      map[int64]int64
	uploads          map[int64]*api.Upload
	passCards        map[int64]*api.PassCardInfo
	projectPassCards map[int64]int64
}

// InsertPass ...
func (m *Memory) InsertPass(ctx context.Context, serialNumber, authToken, passTypeIdentifier string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	pass := &pass{
		serialNumber,
		authToken,
		passTypeIdentifier,
		time.Now(),
	}
	m.passes[pass.serial] = pass
	return nil
}

// UpdatePass ...
func (m *Memory) UpdatePass(ctx context.Context, serialNumber string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	pass, ok := m.passes[serialNumber]
	if !ok {
		return fmt.Errorf("pass %q not found", serialNumber)
	}
	pass.updated = time.Now()
	m.passes[serialNumber] = pass
	return nil
}

// FindPass ...
func (m *Memory) FindPass(ctx context.Context, serialNumber, authToken, passTypeIdentifier string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.passes[serialNumber]
	if !ok {
		return false, fmt.Errorf("pass %q not found", serialNumber)
	}
	return true, nil
}

// FindRegistration ...
func (m *Memory) FindRegistration(ctx context.Context, deviceID, serialNumber string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.regs[deviceID]
	if !ok {
		return false, nil
	}
	return true, nil
}

// FindRegistrationBySerialNumber ...
func (m *Memory) FindRegistrationBySerialNumber(ctx context.Context, serialNumber string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, reg := range m.regs {
		if reg.serial == serialNumber {
			return true, nil
		}
	}
	return false, nil
}

// FindPushToken ...
func (m *Memory) FindPushToken(ctx context.Context, serialNumber string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, reg := range m.regs {
		if reg.serial == serialNumber {
			return reg.push, nil
		}
	}
	return "", store.ErrNotFound
}

// FindSerialNumbers ...
func (m *Memory) FindSerialNumbers(ctx context.Context, deviceID, passTypeIdentifier, tag string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	reg, ok := m.regs[deviceID]
	if !ok {
		return nil, fmt.Errorf("registration %q not found", deviceID)
	}
	return []string{reg.serial}, nil
}

// LatestPass ...
func (m *Memory) LatestPass(ctx context.Context, serialNumber, authToken, passTypeIdentifier string) (time.Time, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pass, ok := m.passes[serialNumber]
	if !ok {
		return time.Time{}, fmt.Errorf("pass %q not found", serialNumber)
	}
	return pass.updated, nil
}

// InsertRegistration ...
func (m *Memory) InsertRegistration(ctx context.Context, deviceID, pushToken, serialNumber, passTypeIdentifier string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	reg := &reg{
		serialNumber,
		deviceID,
		pushToken,
		passTypeIdentifier,
	}
	m.regs[deviceID] = reg
	return nil
}

// DeleteRegistration ...
func (m *Memory) DeleteRegistration(ctx context.Context, deviceID, serialNumber, passTypeIdentifier string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.regs, deviceID)
	return true, nil
}

// InsertLog ...
func (m *Memory) InsertLog(ctx context.Context, remoteAddr, requestID, message string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil
}

// IsUsernameExists ...
func (m *Memory) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, user := range m.users {
		if user.Username == username {
			return true, nil
		}
	}

	return false, nil
}

// IsEmailExists ...
func (m *Memory) IsEmailExists(ctx context.Context, email string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, user := range m.users {
		if user.Email == email {
			return true, nil
		}
	}

	return false, nil
}

// SaveNewUser ...
func (m *Memory) SaveNewUser(ctx context.Context, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.ID]; exists {
		return errors.New("user exists")
	}

	m.users[user.ID] = user

	return nil
}

// LoadUser ...
func (m *Memory) LoadUser(ctx context.Context, id int64) (*api.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not exists")
	}

	return u, nil
}

// LoadUserByUsernameOrEmail ...
func (m *Memory) LoadUserByUsernameOrEmail(ctx context.Context, input string) (*api.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.users {
		if u.Username == input {
			return u, nil
		}
	}

	for _, u := range m.users {
		if u.Email == input {
			return u, nil
		}
	}

	return nil, store.ErrNotFound
}

// LoadUserByConfirmationToken ...
func (m *Memory) LoadUserByConfirmationToken(ctx context.Context, token string) (*api.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.users {
		if u.ConfirmationToken == token {
			return u, nil
		}
	}

	return nil, store.ErrNotFound
}

// LoadUserByRecoveryToken ...
func (m *Memory) LoadUserByRecoveryToken(ctx context.Context, token string) (*api.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.users {
		if u.RecoveryToken == token {
			return u, nil
		}
	}

	return nil, store.ErrNotFound
}

// LoadUserByEmailChangeToken ...
func (m *Memory) LoadUserByEmailChangeToken(ctx context.Context, token string) (*api.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.users {
		if u.EmailChangeToken == token {
			return u, nil
		}
	}

	return nil, store.ErrNotFound
}

// Authenticate ...
func (m *Memory) Authenticate(ctx context.Context, password string, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not exists")
	}

	if ok := user.CheckPassword(password); !ok {
		return store.ErrWrongPassword
	}

	now := time.Now()
	user.LastSignInAt = &now
	m.users[user.ID] = user

	return nil
}

// ConfirmUser ...
func (m *Memory) ConfirmUser(ctx context.Context, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	user.ConfirmedAt = &now
	user.ConfirmationToken = ""
	m.users[user.ID] = user

	return nil
}

// SetConfirmationToken ...
func (m *Memory) SetConfirmationToken(ctx context.Context, confirm api.Confirmation, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	token := secure.Token()

	if confirm == api.SignUpConfirmation {
		user.ConfirmationSentAt = &now
	} else if confirm == api.InviteConfirmation {
		user.InvitedAt = &now
	}

	user.ConfirmationToken = token
	m.users[user.ID] = user

	return nil
}

// RecoverUser ...
func (m *Memory) RecoverUser(ctx context.Context, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.RecoveryToken = ""
	m.users[user.ID] = user

	return nil
}

// SetRecoveryToken ...
func (m *Memory) SetRecoveryToken(ctx context.Context, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	token := secure.Token()

	user.RecoverySentAt = &now
	user.RecoveryToken = token
	m.users[user.ID] = user

	return nil
}

// ConfirmEmailChange ...
func (m *Memory) ConfirmEmailChange(ctx context.Context, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.Email = user.EmailChange
	user.EmailChange = ""
	user.EmailChangeToken = ""
	m.users[user.ID] = user

	return nil
}

// SetEmailChangeToken ...
func (m *Memory) SetEmailChangeToken(ctx context.Context, email string, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	token := secure.Token()

	user.EmailChangeSentAt = &now
	user.EmailChange = email
	user.EmailChangeToken = token
	m.users[user.ID] = user

	return nil
}

// UpdateUsername ...
func (m *Memory) UpdateUsername(ctx context.Context, username string, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.Username = username
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user

	return nil
}

// UpdatePassword ...
func (m *Memory) UpdatePassword(ctx context.Context, hash string, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.PasswordHash = hash
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user

	return nil
}

// UpdateUserMetaData ...
func (m *Memory) UpdateUserMetaData(ctx context.Context, data map[string]interface{}, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.UserMetaData = data
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user

	return nil
}

// UpdateAppMetaData ...
func (m *Memory) UpdateAppMetaData(ctx context.Context, data map[string]interface{}, user *api.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.AppMetaData = data
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user

	return nil
}
