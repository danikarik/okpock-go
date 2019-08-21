package sequel_test

import (
	"context"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store/sequel"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsOrganizationExists(t *testing.T) {
	testCases := []struct {
		Name      string
		SavedOrgs []string
		NewOrg    string
		Expected  bool
	}{
		{
			Name:      "NotTaken",
			SavedOrgs: []string{"GreatApp"},
			NewOrg:    "SuperApp",
			Expected:  false,
		},
		{
			Name:      "Taken",
			SavedOrgs: []string{"CouponApp"},
			NewOrg:    "CouponApp",
			Expected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable, tempOrganizationsTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u, err := api.NewUser(fakeUsername(), fakeEmail(), "test", nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			for _, title := range tc.SavedOrgs {
				org, err := api.NewOrganization(u.ID, title, uuid.NewV4().String(), nil)
				if !assert.NoError(err) {
					return
				}

				err = db.SaveNewOrganization(ctx, org)
				if !assert.NoError(err) {
					return
				}
			}

			exists, err := db.IsOrganizationExists(ctx, u.ID, tc.NewOrg)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestSaveNewOrganization(t *testing.T) {
	testCases := []struct {
		Name      string
		SavedOrgs []string
		NewOrg    string
	}{
		{
			Name:      "NewOrganization",
			SavedOrgs: []string{},
			NewOrg:    "SuperApp",
		},
		{
			Name:      "WithExistingOrganizations",
			SavedOrgs: []string{"CouponApp"},
			NewOrg:    "BoardingPassApp",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable, tempOrganizationsTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u, err := api.NewUser(fakeUsername(), fakeEmail(), "test", nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			for _, title := range tc.SavedOrgs {
				org, err := api.NewOrganization(u.ID, title, uuid.NewV4().String(), nil)
				if !assert.NoError(err) {
					return
				}

				err = db.SaveNewOrganization(ctx, org)
				if !assert.NoError(err) {
					return
				}
			}

			org, err := api.NewOrganization(u.ID, tc.NewOrg, uuid.NewV4().String(), nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewOrganization(ctx, org)
			if !assert.NoError(err) {
				return
			}

			loaded, err := db.LoadOrganization(ctx, org.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(org.ID, loaded.ID)
			assert.Equal(org.Title, loaded.Title)
			assert.Equal(org.Description, loaded.Description)

			loadedOrgs, err := db.LoadOrganizations(ctx, u.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Len(loadedOrgs, len(tc.SavedOrgs)+1)
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	testCases := []struct {
		Name    string
		Org     string
		NewDesc string
		NewData map[string]interface{}
	}{
		{
			Name:    "Coupon",
			Org:     "CouponOrg",
			NewDesc: uuid.NewV4().String(),
			NewData: map[string]interface{}{"quota": 100},
		},
		{
			Name:    "Ticket",
			Org:     "EventTicketOrg",
			NewDesc: uuid.NewV4().String(),
			NewData: map[string]interface{}{"quota": 200},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable, tempOrganizationsTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u, err := api.NewUser(fakeUsername(), fakeEmail(), "test", nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			org, err := api.NewOrganization(u.ID, tc.Org, uuid.NewV4().String(), nil)
			if !assert.NoError(err) {
				return
			}

			err = db.SaveNewOrganization(ctx, org)
			if !assert.NoError(err) {
				return
			}

			err = db.UpdateOrganizationDescription(ctx, tc.NewDesc, org)
			if !assert.NoError(err) {
				return
			}

			err = db.UpdateOrganizationMetaData(ctx, tc.NewData, org)
			if !assert.NoError(err) {
				return
			}

			loaded, err := db.LoadOrganization(ctx, org.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(org.ID, loaded.ID)
			assert.Equal(org.Title, loaded.Title)
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
