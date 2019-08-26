package sequel_test

import (
	"context"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store/sequel"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsProjectExists(t *testing.T) {
	type project struct {
		Desc string
		Type api.PassType
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
				Desc: "Free Coupon",
				Type: api.Coupon,
			},
			Requested: project{
				Desc: "Boarding Pass",
				Type: api.BoardingPass,
			},
			Expected: false,
		},
		{
			Name: "Taken",
			Existing: project{
				Desc: "Free Auction",
				Type: api.Coupon,
			},
			Requested: project{
				Desc: "Free Auction",
				Type: api.Coupon,
			},
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable, tempOrganizationsTable, tempProjectsTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u := api.NewUser(fakeUsername(), fakeEmail(), "test", nil)

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			project := api.NewProject(org.ID, tc.Existing.Desc, tc.Existing.Type)

			err = db.SaveNewProject(ctx, project)
			if !assert.NoError(err) {
				return
			}

			exists, err := db.IsProjectExists(ctx, org.ID, tc.Requested.Desc, tc.Requested.Type)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestSaveNewProject(t *testing.T) {
	type project struct {
		Desc string
		Type api.PassType
	}

	testCases := []struct {
		Name          string
		NewProject    project
		SavedProjects []project
	}{
		{
			Name: "NoExistingProjects",
			NewProject: project{
				Desc: "Free Coupon",
				Type: api.Coupon,
			},
			SavedProjects: []project{},
		},
		{
			Name: "WithExistingProjects",
			NewProject: project{
				Desc: "Boarding Pass",
				Type: api.BoardingPass,
			},
			SavedProjects: []project{
				project{
					Desc: "Generic",
					Type: api.Generic,
				},
				project{
					Desc: "Event",
					Type: api.EventTicket,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable, tempOrganizationsTable, tempProjectsTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u := api.NewUser(fakeUsername(), fakeEmail(), "test", nil)

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			org := api.NewOrganization(u.ID, uuid.NewV4().String(), uuid.NewV4().String(), nil)

			err = db.SaveNewOrganization(ctx, org)
			if !assert.NoError(err) {
				return
			}

			for _, project := range tc.SavedProjects {
				p := api.NewProject(org.ID, project.Desc, project.Type)

				err = db.SaveNewProject(ctx, p)
				if !assert.NoError(err) {
					return
				}
			}

			p := api.NewProject(org.ID, tc.NewProject.Desc, tc.NewProject.Type)

			err = db.SaveNewProject(ctx, p)
			if !assert.NoError(err) {
				return
			}

			loaded, err := db.LoadProject(ctx, p.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(p.ID, loaded.ID)
			assert.Equal(p.Description, loaded.Description)
			assert.Equal(p.PassType, loaded.PassType)

			loadedProjects, err := db.LoadProjects(ctx, u.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Len(loadedProjects, len(tc.SavedProjects)+1)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	type project struct {
		Desc string
		Type api.PassType
	}

	testCases := []struct {
		Name    string
		Project project
		NewDesc string
	}{
		{
			Name: "Coupon",
			Project: project{
				Desc: "Free Coupon",
				Type: api.Coupon,
			},
			NewDesc: "Free Auction",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable, tempOrganizationsTable, tempProjectsTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u := api.NewUser(fakeUsername(), fakeEmail(), "test", nil)

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			org := api.NewOrganization(u.ID, uuid.NewV4().String(), uuid.NewV4().String(), nil)

			err = db.SaveNewOrganization(ctx, org)
			if !assert.NoError(err) {
				return
			}

			p := api.NewProject(org.ID, tc.Project.Desc, tc.Project.Type)

			err = db.SaveNewProject(ctx, p)
			if !assert.NoError(err) {
				return
			}

			err = db.UpdateProjectDescription(ctx, tc.NewDesc, p)
			if !assert.NoError(err) {
				return
			}

			loaded, err := db.LoadProject(ctx, p.ID)
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
		Desc string
		Type api.PassType
	}

	testCases := []struct {
		Name    string
		Project project
		NewKey  string
	}{
		{
			Name: "Coupon",
			Project: project{
				Desc: "Free Coupon",
				Type: api.Coupon,
			},
			NewKey: uuid.NewV4().String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			schema := []string{tempUsersTable, tempOrganizationsTable, tempProjectsTable}
			data := []string{}

			conn, err := executeTempScripts(ctx, t, schema, data)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			u := api.NewUser(fakeUsername(), fakeEmail(), "test", nil)

			err = db.SaveNewUser(ctx, u)
			if !assert.NoError(err) {
				return
			}

			org := api.NewOrganization(u.ID, uuid.NewV4().String(), uuid.NewV4().String(), nil)

			err = db.SaveNewOrganization(ctx, org)
			if !assert.NoError(err) {
				return
			}

			p := api.NewProject(org.ID, tc.Project.Desc, tc.Project.Type)

			err = db.SaveNewProject(ctx, p)
			if !assert.NoError(err) {
				return
			}

			err = db.SetBackgroundImage(ctx, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = db.SetFooterImage(ctx, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = db.SetIconImage(ctx, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = db.SetStripImage(ctx, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			loaded, err := db.LoadProject(ctx, p.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(p.ID, loaded.ID)
			assert.Equal(tc.NewKey, loaded.BackgroundImage)
			assert.Equal(tc.NewKey, loaded.FooterImage)
			assert.Equal(tc.NewKey, loaded.IconImage)
			assert.Equal(tc.NewKey, loaded.StripImage)
		})
	}
}
