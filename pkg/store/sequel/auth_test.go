package sequel_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
	"github.com/danikarik/okpock/pkg/store/sequel"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func fakeUsername() string {
	return uuid.NewV4().String()
}

func fakeEmail() string {
	return uuid.NewV4().String() + "@example.com"
}

func TestUsernameExists(t *testing.T) {
	testCases := []struct {
		Name           string
		SavedUsernames []string
		NewUsername    string
		Expected       bool
	}{
		{
			Name:           "NotTaken",
			SavedUsernames: []string{"mufasa"},
			NewUsername:    "simba",
			Expected:       false,
		},
		{
			Name:           "Taken",
			SavedUsernames: []string{"mufasa", "simba"},
			NewUsername:    "simba",
			Expected:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable}
			data := []string{}

			for _, uname := range tc.SavedUsernames {
				sql := fmt.Sprintf(
					insertUsersTable,
					uuid.NewV4().String()+"@example.com",
					uname,
					"test",
				)
				data = append(data, sql)
			}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			exists, err := db.IsUsernameExists(ctx, tc.NewUsername)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestEmailExists(t *testing.T) {
	testCases := []struct {
		Name        string
		SavedEmails []string
		NewEmail    string
		Expected    bool
	}{
		{
			Name:        "NotTaken",
			SavedEmails: []string{"mufasa@example.com"},
			NewEmail:    "simba@example.com",
			Expected:    false,
		},
		{
			Name:        "Taken",
			SavedEmails: []string{"mufasa@example.com", "simba@example.com"},
			NewEmail:    "simba@example.com",
			Expected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable}
			data := []string{}

			for _, email := range tc.SavedEmails {
				sql := fmt.Sprintf(
					insertUsersTable,
					email,
					uuid.NewV4().String(),
					"test",
				)
				data = append(data, sql)
			}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			exists, err := db.IsEmailExists(ctx, tc.NewEmail)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestSaveNewUser(t *testing.T) {
	type user struct {
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
					Username: "mufasa",
					Email:    "mufasa@jungle.com",
					Password: "king",
				},
				user{
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
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			for _, user := range tc.SavedUsers {
				u, err := api.NewUser(user.Username, user.Email, user.Password, nil)
				if !assert.NoError(err) {
					return
				}

				err = db.SaveNewUser(ctx, u)
				if !assert.NoError(err) {
					return
				}
				if !assert.True(u.ID > 0) {
					return
				}
			}

			exists, err := db.IsEmailExists(ctx, tc.Key)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestSaveInvalidUser(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	schema := []string{tempUsersTable}
	data := []string{}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	user := &api.User{}
	err = db.SaveNewUser(ctx, user)
	assert.Error(err)
}

func TestLoadUser(t *testing.T) {
	testCases := []struct {
		Name string
	}{
		{
			Name: "LoadUser",
		},
		{
			Name: "LoadUserByUsername",
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
		{
			Name: "LoadUserByEmailChangeToken",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			val := uuid.NewV4().String()
			userData := map[string]interface{}{"key": val}
			u, err := api.NewUser(
				fakeUsername(),
				fakeEmail(),
				"test",
				userData,
			)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			err = db.ConfirmUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			var loaded *api.User
			{
				switch tc.Name {
				case "LoadUser":
					loaded, err = db.LoadUser(ctx, u.ID)
					break
				case "LoadUserByUsername":
					loaded, err = db.LoadUserByUsernameOrEmail(ctx, u.Username)
					break
				case "LoadUserByEmail":
					loaded, err = db.LoadUserByUsernameOrEmail(ctx, u.Email)
					break
				case "LoadUserByConfirmationToken":
					err = db.SetConfirmationToken(ctx, api.SignUpConfirmation, u)
					if !assert.NoError(err) {
						return
					}
					loaded, err = db.LoadUserByConfirmationToken(ctx, u.GetConfirmationToken())
					break
				case "LoadUserByRecoveryToken":
					err = db.SetRecoveryToken(ctx, u)
					if !assert.NoError(err) {
						return
					}
					loaded, err = db.LoadUserByRecoveryToken(ctx, u.GetRecoveryToken())
					break
				case "LoadUserByEmailChangeToken":
					err = db.SetEmailChangeToken(ctx, fakeEmail(), u)
					if !assert.NoError(err) {
						return
					}
					loaded, err = db.LoadUserByEmailChangeToken(ctx, u.GetEmailChangeToken())
					break
				default:
					err = store.ErrNotFound
					break
				}
			}

			if assert.NoError(err) {
				assert.Equal(u.ID, loaded.ID)
				assert.Equal(u.Username, loaded.Username)
				assert.Equal(u.Email, loaded.Email)
				assert.Equal(val, loaded.UserMetaData["key"])
				assert.True(loaded.IsConfirmed())
				assert.False(loaded.CreatedAt.IsZero())
				assert.False(loaded.UpdatedAt.IsZero())
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	type user struct {
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
				Username: fakeUsername(),
				Email:    fakeEmail(),
				Password: "test",
			},
			InputPassword: "test",
			HasError:      false,
		},
		{
			Name: "WrongPassword",
			User: user{
				Username: fakeUsername(),
				Email:    fakeEmail(),
				Password: "test",
			},
			InputPassword: "test2",
			HasError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)
			now := time.Now()

			schema := []string{tempUsersTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			err = db.Authenticate(ctx, tc.InputPassword, u)
			if tc.HasError {
				assert.Error(err)
			} else {
				if assert.NoError(err) {
					assert.False(u.LastSignInAt.IsZero())
					assert.True(u.LastSignInAt.Sub(now) > 0)

					loaded, err := db.LoadUser(ctx, u.ID)
					if !assert.NoError(err) {
						return
					}
					assert.False(loaded.LastSignInAt.IsZero())
					assert.True(loaded.LastSignInAt.Add(1*time.Second).Sub(now) > 0)
				}
			}
		})
	}
}

func TestConfirmUser(t *testing.T) {
	type user struct {
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
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			if tc.Confirm {
				err = db.ConfirmUser(ctx, u)
				if !assert.NoError(err) {
					return
				}
			}

			loaded, err := db.LoadUser(ctx, u.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(tc.Expected, loaded.IsConfirmed())
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
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			err = db.SetConfirmationToken(ctx, tc.Confirm, u)
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

			loaded, err := db.LoadUser(ctx, u.ID)
			if !assert.NoError(err) {
				return
			}

			assert.NotEmpty(loaded.GetConfirmationToken())
			assert.Equal(u.GetConfirmationToken(), loaded.GetConfirmationToken())

			if tc.Confirm == api.SignUpConfirmation {
				assert.NotNil(loaded.ConfirmationSentAt)
			} else if tc.Confirm == api.InviteConfirmation {
				assert.NotNil(loaded.InvitedAt)
			}
		})
	}
}

func TestRecoverUser(t *testing.T) {
	type user struct {
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
				Username: "recovered",
				Email:    "recovered@example.com",
				Password: "test",
			},
			Recover: true,
		},
		{
			Name: "NotRecovered",
			User: user{
				Username: "notrecovered",
				Email:    "notrecovered@example.com",
				Password: "test",
			},
			Recover: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			err = db.SetRecoveryToken(ctx, u)
			if !assert.NoError(err) {
				return
			}

			if !assert.NotEmpty(u.GetRecoveryToken()) {
				return
			}

			if tc.Recover {
				err = db.RecoverUser(ctx, u)
				assert.NoError(err)
				assert.Empty(u.GetRecoveryToken())
			} else {
				assert.NotEmpty(u.GetRecoveryToken())
				assert.NotNil(u.RecoverySentAt)
			}

			loaded, err := db.LoadUser(ctx, u.ID)
			if !assert.NoError(err) {
				return
			}

			if tc.Recover {
				assert.Empty(loaded.GetRecoveryToken())
			} else {
				assert.NotEmpty(loaded.GetRecoveryToken())
				assert.NotNil(loaded.RecoverySentAt)
			}
		})
	}
}

func TestEmailChange(t *testing.T) {
	type user struct {
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
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			err = db.SetEmailChangeToken(ctx, tc.User.NewEmail, u)
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
				err = db.ConfirmEmailChange(ctx, u)
				assert.NoError(err)
				assert.Equal(tc.User.NewEmail, u.Email)
			} else {
				assert.NotEmpty(u.GetEmailChangeToken())
				assert.NotNil(u.EmailChangeSentAt)
			}

			loaded, err := db.LoadUser(ctx, u.ID)
			if !assert.NoError(err) {
				return
			}

			if tc.Confirm {
				assert.Equal(tc.User.NewEmail, loaded.Email)
			} else {
				assert.NotEmpty(loaded.GetEmailChangeToken())
				assert.NotNil(loaded.EmailChangeSentAt)
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	schema := []string{tempUsersTable}
	data := []string{}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	user := struct {
		Username    string
		Email       string
		Password    string
		NewPassword string
	}{
		Username:    "passwordchange",
		Email:       "passwordchange@example.com",
		Password:    "test",
		NewPassword: "newpass",
	}

	u, err := api.NewUser(user.Username, user.Email, user.Password, nil)
	if !assert.NoError(err) {
		return
	}

	err = db.SaveNewUser(ctx, u)
	if !assert.NoError(err) {
		return
	}

	err = db.UpdatePassword(ctx, user.NewPassword, u)
	if !assert.NoError(err) {
		return
	}

	loaded, err := db.LoadUser(ctx, u.ID)
	if !assert.NoError(err) {
		return
	}

	ok := loaded.CheckPassword(user.NewPassword)
	assert.True(ok)
}

func TestUpdateMetaData(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	schema := []string{tempUsersTable}
	data := []string{}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	user := struct {
		Username    string
		Email       string
		Password    string
		UserDataKey string
		AppDataKey  string
	}{
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

	err = db.SaveNewUser(ctx, u)
	if !assert.NoError(err) {
		return
	}

	userData := map[string]interface{}{user.UserDataKey: user.UserDataKey}
	err = db.UpdateUserMetaData(ctx, userData, u)
	if !assert.NoError(err) {
		return
	}

	appData := map[string]interface{}{user.AppDataKey: user.AppDataKey}
	err = db.UpdateAppMetaData(ctx, appData, u)
	if !assert.NoError(err) {
		return
	}

	loaded, err := db.LoadUser(ctx, u.ID)
	if !assert.NoError(err) {
		return
	}

	v, ok := loaded.UserMetaData[user.UserDataKey]
	if !assert.True(ok) {
		return
	}
	if !assert.Equal(user.UserDataKey, v) {
		return
	}

	v, ok = loaded.AppMetaData[user.AppDataKey]
	if !assert.True(ok) {
		return
	}
	if !assert.Equal(user.AppDataKey, v) {
		return
	}
}
