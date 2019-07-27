package memory_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
	"github.com/danikarik/okpock/pkg/store/memory"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertPass(t *testing.T) {
	var (
		ctx                = context.Background()
		mock               = memory.New()
		serialNumber       = uuid.NewV4().String()
		authToken          = uuid.NewV4().String()
		passTypeIdentifier = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertPass(ctx, serialNumber, authToken, passTypeIdentifier)
	assert.NoError(err)
}

func TestUpdatePass(t *testing.T) {
	var (
		ctx                = context.Background()
		mock               = memory.New()
		serialNumber       = uuid.NewV4().String()
		authToken          = uuid.NewV4().String()
		passTypeIdentifier = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertPass(ctx, serialNumber, authToken, passTypeIdentifier)
	assert.NoError(err)
	err = mock.UpdatePass(ctx, serialNumber)
	assert.NoError(err)
}

func TestFindPass(t *testing.T) {
	var (
		ctx                = context.Background()
		mock               = memory.New()
		serialNumber       = uuid.NewV4().String()
		authToken          = uuid.NewV4().String()
		passTypeIdentifier = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertPass(ctx, serialNumber, authToken, passTypeIdentifier)
	assert.NoError(err)
	res, err := mock.FindPass(ctx, serialNumber, authToken, passTypeIdentifier)
	assert.NoError(err)
	assert.True(res)
}

func TestFindRegistration(t *testing.T) {
	var (
		ctx                = context.Background()
		mock               = memory.New()
		deviceID           = uuid.NewV4().String()
		serialNumber       = uuid.NewV4().String()
		pushToken          = uuid.NewV4().String()
		passTypeIdentifier = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertRegistration(ctx, deviceID, pushToken, serialNumber, passTypeIdentifier)
	assert.NoError(err)
	res, err := mock.FindRegistration(ctx, deviceID, serialNumber)
	assert.NoError(err)
	assert.True(res)
}

func TestFindSerialNumbers(t *testing.T) {
	var (
		ctx                = context.Background()
		mock               = memory.New()
		deviceID           = uuid.NewV4().String()
		serialNumber       = uuid.NewV4().String()
		pushToken          = uuid.NewV4().String()
		passTypeIdentifier = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertRegistration(ctx, deviceID, pushToken, serialNumber, passTypeIdentifier)
	assert.NoError(err)
	serials, err := mock.FindSerialNumbers(ctx, deviceID, passTypeIdentifier, "")
	assert.NoError(err)
	assert.NotEmpty(serials)
}

func TestLatestPass(t *testing.T) {
	var (
		ctx                = context.Background()
		mock               = memory.New()
		serialNumber       = uuid.NewV4().String()
		authToken          = uuid.NewV4().String()
		passTypeIdentifier = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertPass(ctx, serialNumber, authToken, passTypeIdentifier)
	assert.NoError(err)
	ts, err := mock.LatestPass(ctx, serialNumber, authToken, passTypeIdentifier)
	assert.NoError(err)
	assert.NotNil(ts)
}

func TestInsertRegistration(t *testing.T) {
	var (
		ctx                = context.Background()
		mock               = memory.New()
		deviceID           = uuid.NewV4().String()
		serialNumber       = uuid.NewV4().String()
		pushToken          = uuid.NewV4().String()
		passTypeIdentifier = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertRegistration(ctx, deviceID, pushToken, serialNumber, passTypeIdentifier)
	assert.NoError(err)
}

func TestDeleteRegistration(t *testing.T) {
	var (
		ctx                = context.Background()
		mock               = memory.New()
		deviceID           = uuid.NewV4().String()
		serialNumber       = uuid.NewV4().String()
		pushToken          = uuid.NewV4().String()
		passTypeIdentifier = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertRegistration(ctx, deviceID, pushToken, serialNumber, passTypeIdentifier)
	assert.NoError(err)
	res, err := mock.DeleteRegistration(ctx, deviceID, serialNumber, passTypeIdentifier)
	assert.NoError(err)
	assert.True(res)
}

func TestInsertLog(t *testing.T) {
	var (
		ctx        = context.Background()
		mock       = memory.New()
		remoteAddr = "remote.host"
		requestID  = uuid.NewV4().String()
		message    = "test"
	)
	assert := assert.New(t)
	err := mock.InsertLog(ctx, remoteAddr, requestID, message)
	assert.NoError(err)
}

func TestUsernameExists(t *testing.T) {
	type user struct {
		ID       int64
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name       string
		SavedUsers []user
		Key        string
		Expected   bool
	}{
		{
			Name: "NotTaken",
			SavedUsers: []user{
				user{
					ID:       1,
					Username: "mufasa",
					Email:    "mufasa@jungle.com",
					Password: "king",
				},
			},
			Key:      "simba",
			Expected: false,
		},
		{
			Name: "Taken",
			SavedUsers: []user{
				user{
					ID:       2,
					Username: "mufasa",
					Email:    "mufasa@jungle.com",
					Password: "king",
				},
				user{
					ID:       3,
					Username: "simba",
					Email:    "simba@jungle.com",
					Password: "prince",
				},
			},
			Key:      "simba",
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			for _, user := range tc.SavedUsers {
				u, err := api.NewUser(user.Username, user.Email, user.Password, nil)
				if !assert.NoError(err) {
					return
				}
				u.ID = user.ID

				err = mock.SaveNewUser(ctx, u)
				if !assert.NoError(err) {
					return
				}
			}

			exists, err := mock.IsUsernameExists(ctx, tc.Key)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestEmailExists(t *testing.T) {
	type user struct {
		ID       int64
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name       string
		SavedUsers []user
		Key        string
		Expected   bool
	}{
		{
			Name: "NotTaken",
			SavedUsers: []user{
				user{
					ID:       4,
					Username: "mufasa",
					Email:    "mufasa@jungle.com",
					Password: "king",
				},
			},
			Key:      "simba@jungle.com",
			Expected: false,
		},
		{
			Name: "Taken",
			SavedUsers: []user{
				user{
					ID:       5,
					Username: "mufasa",
					Email:    "mufasa@jungle.com",
					Password: "king",
				},
				user{
					ID:       6,
					Username: "simba",
					Email:    "simba@jungle.com",
					Password: "prince",
				},
			},
			Key:      "simba@jungle.com",
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			for _, user := range tc.SavedUsers {
				u, err := api.NewUser(user.Username, user.Email, user.Password, nil)
				if !assert.NoError(err) {
					return
				}
				u.ID = user.ID

				err = mock.SaveNewUser(ctx, u)
				if !assert.NoError(err) {
					return
				}
			}

			exists, err := mock.IsEmailExists(ctx, tc.Key)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestSaveNewUser(t *testing.T) {
	type user struct {
		ID       int64
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name       string
		SavedUsers []user
		Key        string
		Expected   bool
	}{
		{
			Name: "NotTaken",
			SavedUsers: []user{
				user{
					ID:       7,
					Username: "mufasa",
					Email:    "mufasa@jungle.com",
					Password: "king",
				},
			},
			Key:      "simba@jungle.com",
			Expected: false,
		},
		{
			Name: "Taken",
			SavedUsers: []user{
				user{
					ID:       8,
					Username: "mufasa",
					Email:    "mufasa@jungle.com",
					Password: "king",
				},
				user{
					ID:       9,
					Username: "simba",
					Email:    "simba@jungle.com",
					Password: "prince",
				},
			},
			Key:      "simba@jungle.com",
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			for _, user := range tc.SavedUsers {
				u, err := api.NewUser(user.Username, user.Email, user.Password, nil)
				if !assert.NoError(err) {
					return
				}
				u.ID = user.ID

				err = mock.SaveNewUser(ctx, u)
				if !assert.NoError(err) {
					return
				}
			}

			exists, err := mock.IsEmailExists(ctx, tc.Key)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestLoadUser(t *testing.T) {
	testCases := []struct {
		Name string
	}{
		{
			Name: "LoadUser",
		},
		{
			Name: "LoadUserByEmail",
		},
		{
			Name: "LoadUserByConfirmationToken",
		},
		{
			Name: "LoadUserByRecoveryToken",
		},
	}

	for i, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx      = context.Background()
				assert   = assert.New(t)
				mock     = memory.New()
				username = strings.ToLower(tc.Name)
				email    = fmt.Sprintf("%s@example.com", strings.ToLower(tc.Name))
			)

			u, err := api.NewUser(username, email, "test", nil)
			if !assert.NoError(err) {
				return
			}
			u.ID = int64(i + 1)

			err = mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			var loaded *api.User
			{
				switch tc.Name {
				case "LoadUser":
					loaded, err = mock.LoadUser(ctx, u.ID)
					break
				case "LoadUserByEmail":
					loaded, err = mock.LoadUserByUsernameOrEmail(ctx, u.Email)
					break
				case "LoadUserByConfirmationToken":
					err = mock.SetConfirmationToken(ctx, api.SignUpConfirmation, u)
					if !assert.NoError(err) {
						return
					}
					loaded, err = mock.LoadUserByConfirmationToken(ctx, u.GetConfirmationToken())
					break
				case "LoadUserByRecoveryToken":
					err = mock.SetRecoveryToken(ctx, u)
					if !assert.NoError(err) {
						return
					}
					loaded, err = mock.LoadUserByRecoveryToken(ctx, u.GetRecoveryToken())
					break
				default:
					err = store.ErrNotFound
					break
				}
			}

			if assert.NoError(err) {
				assert.Equal(u.ID, loaded.ID)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	type user struct {
		ID       int64
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name          string
		User          user
		InputPassword string
		HasError      bool
	}{
		{
			Name: "CorrectPassword",
			User: user{
				ID:       1,
				Username: "correct",
				Email:    "correct@example.com",
				Password: "test",
			},
			InputPassword: "test",
			HasError:      false,
		},
		{
			Name: "WrongPassword",
			User: user{
				ID:       2,
				Username: "wrong",
				Email:    "wrong@example.com",
				Password: "test",
			},
			InputPassword: "test2",
			HasError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
				now    = time.Now()
			)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			err = mock.Authenticate(ctx, tc.InputPassword, u)
			if tc.HasError {
				assert.Error(err)
			} else {
				if assert.NoError(err) {
					assert.False(u.LastSignInAt.IsZero())
					assert.True(u.LastSignInAt.Sub(now) > 0)
				}
			}
		})
	}
}

func TestConfirmUser(t *testing.T) {
	type user struct {
		ID       int64
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     user
		Confirm  bool
		Expected bool
	}{
		{
			Name: "NotConfirmed",
			User: user{
				ID:       1,
				Username: "notconfirmed",
				Email:    "notconfirmed@example.com",
				Password: "test",
			},
			Confirm:  false,
			Expected: false,
		},
		{
			Name: "Confirmed",
			User: user{
				ID:       2,
				Username: "confirmed",
				Email:    "confirmed@example.com",
				Password: "test",
			},
			Confirm:  true,
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			if tc.Confirm {
				err = mock.ConfirmUser(ctx, u)
				if !assert.NoError(err) {
					return
				}
			}

			assert.Equal(tc.Expected, u.IsConfirmed())
		})
	}
}

func TestSetConfirmationToken(t *testing.T) {
	type user struct {
		ID       int64
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name    string
		User    user
		Confirm api.Confirmation
	}{
		{
			Name: "SignUp",
			User: user{
				ID:       1,
				Username: "signup",
				Email:    "signup@example.com",
				Password: "test",
			},
			Confirm: api.SignUpConfirmation,
		},
		{
			Name: "Invite",
			User: user{
				ID:       2,
				Username: "invite",
				Email:    "invite@example.com",
				Password: "test",
			},
			Confirm: api.InviteConfirmation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			err = mock.SetConfirmationToken(ctx, tc.Confirm, u)
			if !assert.NoError(err) {
				return
			}

			if !assert.NotEmpty(u.GetConfirmationToken()) {
				return
			}

			if tc.Confirm == api.SignUpConfirmation {
				assert.NotNil(u.ConfirmationSentAt)
			} else if tc.Confirm == api.InviteConfirmation {
				assert.NotNil(u.InvitedAt)
			}
		})
	}
}

func TestRecoverUser(t *testing.T) {
	type user struct {
		ID       int64
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name    string
		User    user
		Recover bool
	}{
		{
			Name: "Recovered",
			User: user{
				ID:       1,
				Username: "recovered",
				Email:    "recovered@example.com",
				Password: "test",
			},
			Recover: true,
		},
		{
			Name: "NotRecovered",
			User: user{
				ID:       2,
				Username: "notrecovered",
				Email:    "notrecovered@example.com",
				Password: "test",
			},
			Recover: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			err = mock.SetRecoveryToken(ctx, u)
			if !assert.NoError(err) {
				return
			}

			if !assert.NotEmpty(u.GetRecoveryToken()) {
				return
			}

			if tc.Recover {
				err = mock.RecoverUser(ctx, u)
				assert.NoError(err)
				assert.Empty(u.GetRecoveryToken())
			} else {
				assert.NotEmpty(u.GetRecoveryToken())
				assert.NotNil(u.RecoverySentAt)
			}
		})
	}
}

func TestEmailChange(t *testing.T) {
	type user struct {
		ID       int64
		Username string
		Email    string
		NewEmail string
		Password string
	}

	testCases := []struct {
		Name    string
		User    user
		Confirm bool
	}{
		{
			Name: "EmailConfirmed",
			User: user{
				ID:       1,
				Username: "emailconfirmed",
				Email:    "emailconfirmed@example.com",
				NewEmail: "newemailconfirmed@example.com",
				Password: "test",
			},
			Confirm: true,
		},
		{
			Name: "EmailNotConfirmed",
			User: user{
				ID:       2,
				Username: "emailnotconfirmed",
				Email:    "emailnotconfirmed@example.com",
				NewEmail: "newemailnotconfirmed@example.com",
				Password: "test",
			},
			Confirm: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			err = mock.SetEmailChangeToken(ctx, tc.User.NewEmail, u)
			if !assert.NoError(err) {
				return
			}

			if !assert.NotEmpty(u.GetEmailChangeToken()) {
				return
			}
			if !assert.Equal(tc.User.NewEmail, u.GetEmailChange()) {
				return
			}

			if tc.Confirm {
				err = mock.ConfirmEmailChange(ctx, u)
				assert.NoError(err)
				assert.Equal(tc.User.NewEmail, u.Email)
			} else {
				assert.NotEmpty(u.GetEmailChangeToken())
				assert.NotNil(u.EmailChangeSentAt)
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		mock   = memory.New()
	)

	user := struct {
		ID          int64
		Username    string
		Email       string
		Password    string
		NewPassword string
	}{
		ID:          1,
		Username:    "passwordchange",
		Email:       "passwordchange@example.com",
		Password:    "test",
		NewPassword: "newpass",
	}

	u, err := api.NewUser(user.Username, user.Email, user.Password, nil)
	if !assert.NoError(err) {
		return
	}

	err = mock.SaveNewUser(ctx, u)
	if !assert.NoError(err) {
		return
	}
	u.ID = user.ID

	err = mock.UpdatePassword(ctx, user.NewPassword, u)
	if !assert.NoError(err) {
		return
	}

	ok := u.CheckPassword(user.NewPassword)
	assert.True(ok)
}

func TestUpdateMetaData(t *testing.T) {
	var (
		ctx    = context.Background()
		assert = assert.New(t)
		mock   = memory.New()
	)

	user := struct {
		ID          int64
		Username    string
		Email       string
		Password    string
		UserDataKey string
		AppDataKey  string
	}{
		ID:          111,
		Username:    "metadata",
		Email:       "metadata@example.com",
		Password:    "test",
		UserDataKey: "user_id",
		AppDataKey:  "app_version",
	}

	u, err := api.NewUser(user.Username, user.Email, user.Password, nil)
	if !assert.NoError(err) {
		return
	}

	err = mock.SaveNewUser(ctx, u)
	if !assert.NoError(err) {
		return
	}
	u.ID = user.ID

	userData := map[string]interface{}{user.UserDataKey: user.UserDataKey}
	err = mock.UpdateUserMetaData(ctx, userData, u)
	if !assert.NoError(err) {
		return
	}

	v, ok := u.UserMetaData[user.UserDataKey]
	if !assert.True(ok) {
		return
	}
	if !assert.Equal(user.UserDataKey, v) {
		return
	}

	appData := map[string]interface{}{user.AppDataKey: user.AppDataKey}
	err = mock.UpdateAppMetaData(ctx, appData, u)
	if !assert.NoError(err) {
		return
	}

	v, ok = u.AppMetaData[user.AppDataKey]
	if !assert.True(ok) {
		return
	}
	if !assert.Equal(user.AppDataKey, v) {
		return
	}
}
