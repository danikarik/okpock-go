package memory_test

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/danikarik/okpock/pkg/secure"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
	"github.com/danikarik/okpock/pkg/store/memory"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func randomID() int64 { return rand.Int63() }

func TestInsertPass(t *testing.T) {
	var (
		ctx          = context.Background()
		mock         = memory.New()
		serialNumber = uuid.NewV4().String()
		authToken    = uuid.NewV4().String()
		passTypeID   = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertPass(ctx, serialNumber, authToken, passTypeID)
	assert.NoError(err)
}

func TestUpdatePass(t *testing.T) {
	var (
		ctx          = context.Background()
		mock         = memory.New()
		serialNumber = uuid.NewV4().String()
		authToken    = uuid.NewV4().String()
		passTypeID   = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertPass(ctx, serialNumber, authToken, passTypeID)
	assert.NoError(err)
	err = mock.UpdatePass(ctx, serialNumber)
	assert.NoError(err)
}

func TestFindPass(t *testing.T) {
	var (
		ctx          = context.Background()
		mock         = memory.New()
		serialNumber = uuid.NewV4().String()
		authToken    = uuid.NewV4().String()
		passTypeID   = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertPass(ctx, serialNumber, authToken, passTypeID)
	assert.NoError(err)
	res, err := mock.FindPass(ctx, serialNumber, authToken, passTypeID)
	assert.NoError(err)
	assert.True(res)
}

func TestFindRegistration(t *testing.T) {
	var (
		ctx          = context.Background()
		mock         = memory.New()
		deviceID     = uuid.NewV4().String()
		serialNumber = uuid.NewV4().String()
		pushToken    = uuid.NewV4().String()
		passTypeID   = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertRegistration(ctx, deviceID, pushToken, serialNumber, passTypeID)
	assert.NoError(err)
	res, err := mock.FindRegistration(ctx, deviceID, serialNumber)
	assert.NoError(err)
	assert.True(res)
}

func TestFindSerialNumbers(t *testing.T) {
	var (
		ctx          = context.Background()
		mock         = memory.New()
		deviceID     = uuid.NewV4().String()
		serialNumber = uuid.NewV4().String()
		pushToken    = uuid.NewV4().String()
		passTypeID   = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertRegistration(ctx, deviceID, pushToken, serialNumber, passTypeID)
	assert.NoError(err)
	serials, err := mock.FindSerialNumbers(ctx, deviceID, passTypeID, "")
	assert.NoError(err)
	assert.NotEmpty(serials)
}

func TestLatestPass(t *testing.T) {
	var (
		ctx          = context.Background()
		mock         = memory.New()
		serialNumber = uuid.NewV4().String()
		authToken    = uuid.NewV4().String()
		passTypeID   = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertPass(ctx, serialNumber, authToken, passTypeID)
	assert.NoError(err)
	ts, err := mock.LatestPass(ctx, serialNumber, authToken, passTypeID)
	assert.NoError(err)
	assert.NotNil(ts)
}

func TestInsertRegistration(t *testing.T) {
	var (
		ctx          = context.Background()
		mock         = memory.New()
		deviceID     = uuid.NewV4().String()
		serialNumber = uuid.NewV4().String()
		pushToken    = uuid.NewV4().String()
		passTypeID   = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertRegistration(ctx, deviceID, pushToken, serialNumber, passTypeID)
	assert.NoError(err)
}

func TestDeleteRegistration(t *testing.T) {
	var (
		ctx          = context.Background()
		mock         = memory.New()
		deviceID     = uuid.NewV4().String()
		serialNumber = uuid.NewV4().String()
		pushToken    = uuid.NewV4().String()
		passTypeID   = "test.passkit"
	)
	assert := assert.New(t)
	err := mock.InsertRegistration(ctx, deviceID, pushToken, serialNumber, passTypeID)
	assert.NoError(err)
	res, err := mock.DeleteRegistration(ctx, deviceID, serialNumber, passTypeID)
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
				u := api.NewUser(user.Username, user.Email, user.Password, nil)
				u.ID = user.ID

				err := mock.SaveNewUser(ctx, u)
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
				u := api.NewUser(user.Username, user.Email, user.Password, nil)
				u.ID = user.ID

				err := mock.SaveNewUser(ctx, u)
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
				u := api.NewUser(user.Username, user.Email, user.Password, nil)
				u.ID = user.ID

				err := mock.SaveNewUser(ctx, u)
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
		{
			Name: "LoadUserByEmailChangeToken",
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

			u := api.NewUser(username, email, "test", nil)
			u.ID = int64(i + 1)

			err := mock.SaveNewUser(ctx, u)
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
				case "LoadUserByEmailChangeToken":
					err = mock.SetEmailChangeToken(ctx, "newemail@example.com", u)
					if !assert.NoError(err) {
						return
					}
					loaded, err = mock.LoadUserByEmailChangeToken(ctx, u.GetEmailChangeToken())
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

			hash, err := secure.NewPassword(tc.User.Password)
			if !assert.NoError(err) {
				return
			}

			u := api.NewUser(tc.User.Username, tc.User.Email, hash, nil)

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

			u := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)

			err := mock.SaveNewUser(ctx, u)
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

			u := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)

			err := mock.SaveNewUser(ctx, u)
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

			u := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)

			err := mock.SaveNewUser(ctx, u)
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

			u := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)

			err := mock.SaveNewUser(ctx, u)
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

func TestUpdateUsername(t *testing.T) {
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
		NewUsername string
	}{
		ID:          1,
		Username:    "usernamechange",
		Email:       "usernamechange@example.com",
		Password:    "test",
		NewUsername: "newusername",
	}

	u := api.NewUser(user.Username, user.Email, user.Password, nil)
	u.ID = user.ID

	err := mock.SaveNewUser(ctx, u)
	if !assert.NoError(err) {
		return
	}

	err = mock.UpdateUsername(ctx, user.NewUsername, u)
	if !assert.NoError(err) {
		return
	}

	assert.Equal(user.NewUsername, u.Username)
	assert.False(u.UpdatedAt.IsZero())
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

	u := api.NewUser(user.Username, user.Email, user.Password, nil)
	u.ID = user.ID

	err := mock.SaveNewUser(ctx, u)
	if !assert.NoError(err) {
		return
	}

	err = mock.UpdatePassword(ctx, user.NewPassword, u)
	if !assert.NoError(err) {
		return
	}

	ok := u.CheckPassword(user.NewPassword)
	assert.True(ok)
	assert.False(u.UpdatedAt.IsZero())
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

	u := api.NewUser(user.Username, user.Email, user.Password, nil)
	u.ID = user.ID

	err := mock.SaveNewUser(ctx, u)
	if !assert.NoError(err) {
		return
	}

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
	assert.False(u.UpdatedAt.IsZero())

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
	assert.False(u.UpdatedAt.IsZero())
}

func TestIsOrganizationExists(t *testing.T) {
	type org struct {
		ID     int64
		Title  string
		Desc   string
		UserID int64
	}

	testCases := []struct {
		Name      string
		Existing  org
		Requested org
		Expected  bool
	}{
		{
			Name: "NotTaken",
			Existing: org{
				ID:     1,
				Title:  "GreatApp",
				Desc:   "Sample Organization",
				UserID: 1,
			},
			Requested: org{
				Title:  "AnotherGreatApp",
				UserID: 2,
			},
			Expected: false,
		},
		{
			Name: "TakenTitle",
			Existing: org{
				ID:     2,
				Title:  "GreatApp",
				Desc:   "Sample Organization",
				UserID: 3,
			},
			Requested: org{
				Title:  "GreatApp",
				UserID: 4,
			},
			Expected: false,
		},
		{
			Name: "Exists",
			Existing: org{
				ID:     3,
				Title:  "GreatApp",
				Desc:   "Sample Organization",
				UserID: 5,
			},
			Requested: org{
				Title:  "GreatApp",
				UserID: 5,
			},
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

			org := api.NewOrganization(tc.Existing.UserID, tc.Existing.Title, tc.Existing.Desc, nil)
			org.ID = tc.Existing.ID

			err := mock.SaveNewOrganization(ctx, org)
			if !assert.NoError(err) {
				return
			}

			exists, err := mock.IsOrganizationExists(ctx, tc.Requested.UserID, tc.Requested.Title)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestSaveNewOrganization(t *testing.T) {
	type org struct {
		ID     int64
		Title  string
		Desc   string
		UserID int64
	}

	testCases := []struct {
		Name      string
		NewOrg    org
		SavedOrgs []org
	}{
		{
			Name: "NoExistingOrgs",
			NewOrg: org{
				ID:    randomID(),
				Title: "GreatOrg",
				Desc:  "Sample Org",
			},
			SavedOrgs: []org{},
		},
		{
			Name: "WithExistingOrgs",
			NewOrg: org{
				ID:    randomID(),
				Title: "AnotherGreatOrg",
				Desc:  "Sample Org",
			},
			SavedOrgs: []org{
				org{
					ID:    randomID(),
					Title: "GreatOrg2",
					Desc:  "Sample Org",
				},
				org{
					ID:    randomID(),
					Title: "GreatOrg3",
					Desc:  "Sample Org",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			var (
				id       = randomID()
				username = fmt.Sprintf("user%d", id)
				email    = fmt.Sprintf("user%d@example.com", id)
			)

			u := api.NewUser(username, email, "test", nil)
			u.ID = id

			err := mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			for _, org := range tc.SavedOrgs {
				o := api.NewOrganization(u.ID, org.Title, org.Desc, nil)
				o.ID = org.ID

				err := mock.SaveNewOrganization(ctx, o)
				if !assert.NoError(err) {
					return
				}
			}

			o := api.NewOrganization(u.ID, tc.NewOrg.Title, tc.NewOrg.Desc, nil)
			o.ID = tc.NewOrg.ID

			err = mock.SaveNewOrganization(ctx, o)
			if !assert.NoError(err) {
				return
			}

			loaded, err := mock.LoadOrganization(ctx, o.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(o.ID, loaded.ID)
			assert.Equal(o.Title, loaded.Title)
			assert.Equal(o.Description, loaded.Description)

			loadedOrgs, err := mock.LoadOrganizations(ctx, u.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Len(loadedOrgs, len(tc.SavedOrgs)+1)
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	type org struct {
		ID     int64
		Title  string
		Desc   string
		UserID int64
	}

	testCases := []struct {
		Name    string
		Org     org
		NewDesc string
		NewData map[string]interface{}
	}{
		{
			Name: "GreatOrg",
			Org: org{
				ID:    randomID(),
				Title: "GreatOrg",
				Desc:  "Sample Org",
			},
			NewDesc: "Updated Description",
			NewData: map[string]interface{}{"quota": 100},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			var (
				id       = randomID()
				username = fmt.Sprintf("user%d", id)
				email    = fmt.Sprintf("user%d@example.com", id)
			)

			u := api.NewUser(username, email, "test", nil)
			u.ID = id

			err := mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			o := api.NewOrganization(u.ID, tc.Org.Title, tc.Org.Desc, nil)
			o.ID = tc.Org.ID

			err = mock.SaveNewOrganization(ctx, o)
			if !assert.NoError(err) {
				return
			}

			err = mock.UpdateOrganizationDescription(ctx, tc.NewDesc, o)
			if !assert.NoError(err) {
				return
			}

			err = mock.UpdateOrganizationMetaData(ctx, tc.NewData, o)
			if !assert.NoError(err) {
				return
			}

			loaded, err := mock.LoadOrganization(ctx, o.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(o.ID, loaded.ID)
			assert.Equal(o.Title, loaded.Title)
			assert.Equal(tc.NewDesc, loaded.Description)

			ok := true
			for k := range tc.NewData {
				if _, has := loaded.MetaData[k]; !has {
					ok = false
				}
			}
			assert.True(ok)
		})
	}
}

func TestIsProjectExists(t *testing.T) {
	type project struct {
		ID    int64
		OrgID int64
		Desc  string
		Type  api.PassType
	}

	testCases := []struct {
		Name      string
		Existing  project
		Requested project
		Expected  bool
	}{
		{
			Name: "NotTaken",
			Existing: project{
				ID:    randomID(),
				OrgID: randomID(),
				Desc:  "Free Coupon",
				Type:  api.Coupon,
			},
			Requested: project{
				ID:    randomID(),
				OrgID: randomID(),
				Desc:  "Boarding Pass",
				Type:  api.BoardingPass,
			},
			Expected: false,
		},
		{
			Name: "TakenDescription",
			Existing: project{
				ID:    randomID(),
				OrgID: randomID(),
				Desc:  "Free Auction",
				Type:  api.Coupon,
			},
			Requested: project{
				ID:    randomID(),
				OrgID: randomID(),
				Desc:  "Free Auction",
				Type:  api.EventTicket,
			},
			Expected: false,
		},
		{
			Name: "Exists",
			Existing: project{
				ID:    randomID(),
				OrgID: 11,
				Desc:  "Free Auction",
				Type:  api.Coupon,
			},
			Requested: project{
				ID:    randomID(),
				OrgID: 11,
				Desc:  "Free Auction",
				Type:  api.Coupon,
			},
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

			project := api.NewProject(tc.Existing.OrgID, tc.Existing.Desc, tc.Existing.Type)
			project.ID = tc.Existing.ID

			err := mock.SaveNewProject(ctx, project)
			if !assert.NoError(err) {
				return
			}

			exists, err := mock.IsProjectExists(ctx, tc.Requested.OrgID, tc.Requested.Desc, tc.Requested.Type)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestSaveNewProject(t *testing.T) {
	type project struct {
		ID    int64
		OrgID int64
		Desc  string
		Type  api.PassType
	}

	testCases := []struct {
		Name          string
		NewProject    project
		SavedProjects []project
	}{
		{
			Name: "NoExistingProjects",
			NewProject: project{
				ID:    randomID(),
				OrgID: randomID(),
				Desc:  "Free Coupon",
				Type:  api.Coupon,
			},
			SavedProjects: []project{},
		},
		{
			Name: "WithExistingProjects",
			NewProject: project{
				ID:    randomID(),
				OrgID: randomID(),
				Desc:  "Boarding Pass",
				Type:  api.BoardingPass,
			},
			SavedProjects: []project{
				project{
					ID:    randomID(),
					OrgID: randomID(),
					Desc:  "Generic",
					Type:  api.Generic,
				},
				project{
					ID:    randomID(),
					OrgID: randomID(),
					Desc:  "Event",
					Type:  api.EventTicket,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			var (
				id       = randomID()
				username = fmt.Sprintf("user%d", id)
				email    = fmt.Sprintf("user%d@example.com", id)
				orgTitle = fmt.Sprintf("title%d", id)
				orgDesc  = fmt.Sprintf("desc%d", id)
			)

			u := api.NewUser(username, email, "test", nil)
			u.ID = id

			err := mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			o := api.NewOrganization(u.ID, orgTitle, orgDesc, nil)
			o.ID = randomID()

			err = mock.SaveNewOrganization(ctx, o)
			if !assert.NoError(err) {
				return
			}

			for _, project := range tc.SavedProjects {
				p := api.NewProject(o.ID, project.Desc, project.Type)
				p.ID = project.ID

				err = mock.SaveNewProject(ctx, p)
				if !assert.NoError(err) {
					return
				}
			}

			p := api.NewProject(o.ID, tc.NewProject.Desc, tc.NewProject.Type)
			p.ID = tc.NewProject.ID

			err = mock.SaveNewProject(ctx, p)
			if !assert.NoError(err) {
				return
			}

			loaded, err := mock.LoadProject(ctx, p.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(p.ID, loaded.ID)
			assert.Equal(p.Description, loaded.Description)
			assert.Equal(p.PassType, loaded.PassType)

			loadedProjects, err := mock.LoadProjects(ctx, u.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Len(loadedProjects, len(tc.SavedProjects)+1)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	type project struct {
		ID    int64
		OrgID int64
		Desc  string
		Type  api.PassType
	}

	testCases := []struct {
		Name    string
		Project project
		NewDesc string
	}{
		{
			Name: "Coupon",
			Project: project{
				ID:    randomID(),
				OrgID: randomID(),
				Desc:  "Free Coupon",
				Type:  api.Coupon,
			},
			NewDesc: "Free Auction",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			var (
				id       = randomID()
				username = fmt.Sprintf("user%d", id)
				email    = fmt.Sprintf("user%d@example.com", id)
				orgTitle = fmt.Sprintf("title%d", id)
				orgDesc  = fmt.Sprintf("desc%d", id)
			)

			u := api.NewUser(username, email, "test", nil)
			u.ID = id

			err := mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			o := api.NewOrganization(u.ID, orgTitle, orgDesc, nil)
			o.ID = randomID()

			err = mock.SaveNewOrganization(ctx, o)
			if !assert.NoError(err) {
				return
			}

			p := api.NewProject(o.ID, tc.Project.Desc, tc.Project.Type)
			p.ID = tc.Project.ID

			err = mock.SaveNewProject(ctx, p)
			if !assert.NoError(err) {
				return
			}

			err = mock.UpdateProjectDescription(ctx, tc.NewDesc, p)
			if !assert.NoError(err) {
				return
			}

			loaded, err := mock.LoadProject(ctx, p.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(p.ID, loaded.ID)
			assert.Equal(p.Description, loaded.Description)
		})
	}
}

func TestSetImage(t *testing.T) {
	type project struct {
		ID    int64
		OrgID int64
		Desc  string
		Type  api.PassType
	}

	testCases := []struct {
		Name    string
		Project project
		NewKey  string
	}{
		{
			Name: "Coupon",
			Project: project{
				ID:    randomID(),
				OrgID: randomID(),
				Desc:  "Free Coupon",
				Type:  api.Coupon,
			},
			NewKey: uuid.NewV4().String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				assert = assert.New(t)
				mock   = memory.New()
			)

			var (
				id       = randomID()
				username = fmt.Sprintf("user%d", id)
				email    = fmt.Sprintf("user%d@example.com", id)
				orgTitle = fmt.Sprintf("title%d", id)
				orgDesc  = fmt.Sprintf("desc%d", id)
			)

			u := api.NewUser(username, email, "test", nil)
			u.ID = id

			err := mock.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			o := api.NewOrganization(u.ID, orgTitle, orgDesc, nil)
			o.ID = randomID()

			err = mock.SaveNewOrganization(ctx, o)
			if !assert.NoError(err) {
				return
			}

			p := api.NewProject(o.ID, tc.Project.Desc, tc.Project.Type)
			p.ID = tc.Project.ID

			err = mock.SaveNewProject(ctx, p)
			if !assert.NoError(err) {
				return
			}

			err = mock.SetBackgroundImage(ctx, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = mock.SetFooterImage(ctx, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = mock.SetIconImage(ctx, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = mock.SetStripImage(ctx, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			loaded, err := mock.LoadProject(ctx, p.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(p.ID, loaded.ID)
			assert.Equal(tc.NewKey, loaded.GetField(api.BackgroundImage))
			assert.Equal(tc.NewKey, loaded.GetField(api.FooterImage))
			assert.Equal(tc.NewKey, loaded.GetField(api.IconImage))
			assert.Equal(tc.NewKey, loaded.GetField(api.StripImage))
		})
	}
}
