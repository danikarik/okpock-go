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
		Title            string
		OrganizationName string
		Desc             string
		Type             api.PassType
	}

	takenTitle := fakeString()

	testCases := []struct {
		Name      string
		Existing  project
		Requested project
		Expected  bool
	}{
		{
			Name: "NotTaken",
			Existing: project{
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Desc:             fakeString(),
				Type:             api.Coupon,
			},
			Requested: project{
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Desc:             fakeString(),
				Type:             api.BoardingPass,
			},
			Expected: false,
		},
		{
			Name: "Taken",
			Existing: project{
				Title:            takenTitle,
				OrganizationName: takenTitle,
				Desc:             takenTitle,
				Type:             api.BoardingPass,
			},
			Requested: project{
				Title:            takenTitle,
				OrganizationName: takenTitle,
				Desc:             takenTitle,
				Type:             api.BoardingPass,
			},
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			conn, err := testConnection(ctx, t)
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

			project := api.NewProject(tc.Existing.Title, tc.Existing.OrganizationName, tc.Existing.Desc, tc.Existing.Type)

			err = db.SaveNewProject(ctx, u, project)
			if !assert.NoError(err) {
				return
			}

			exists, err := db.IsProjectExists(ctx, tc.Requested.Title, tc.Requested.OrganizationName, tc.Requested.Desc, tc.Requested.Type)
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

			conn, err := testConnection(ctx, t)
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

			orgName := fakeString()

			for _, project := range tc.SavedProjects {
				p := api.NewProject(orgName, fakeString(), project.Desc, project.Type)

				err = db.SaveNewProject(ctx, u, p)
				if !assert.NoError(err) {
					return
				}
			}

			p := api.NewProject(orgName, fakeString(), tc.NewProject.Desc, tc.NewProject.Type)

			err = db.SaveNewProject(ctx, u, p)
			if !assert.NoError(err) {
				return
			}

			loaded, err := db.LoadProject(ctx, u, p.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(p.ID, loaded.ID)
			assert.Equal(p.Description, loaded.Description)
			assert.Equal(p.PassType, loaded.PassType)

			loadedProjects, err := db.LoadProjects(ctx, u, nil)
			if !assert.NoError(err) {
				return
			}

			assert.Len(loadedProjects.Data, len(tc.SavedProjects)+1)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	type project struct {
		Title string
		Name  string
		Desc  string
		Type  api.PassType
	}

	testCases := []struct {
		Name       string
		Project    project
		NewTitle   string
		NewOrgName string
		NewDesc    string
	}{
		{
			Name: "Coupon",
			Project: project{
				Title: fakeString(),
				Name:  fakeString(),
				Desc:  fakeString(),
				Type:  api.Coupon,
			},
			NewTitle:   fakeString(),
			NewOrgName: fakeString(),
			NewDesc:    fakeString(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			conn, err := testConnection(ctx, t)
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

			p := api.NewProject(tc.Project.Title,
				tc.Project.Name,
				tc.Project.Desc,
				tc.Project.Type)

			err = db.SaveNewProject(ctx, u, p)
			if !assert.NoError(err) {
				return
			}

			err = db.UpdateProject(ctx, tc.NewTitle, tc.NewOrgName, tc.NewDesc, p)
			if !assert.NoError(err) {
				return
			}

			loaded, err := db.LoadProject(ctx, u, p.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(p.ID, loaded.ID)
			assert.Equal(p.Title, loaded.Title)
			assert.Equal(p.OrganizationName, loaded.OrganizationName)
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
		Size    api.ImageSize
		Project project
		NewKey  string
	}{
		{
			Name: "Coupon1x",
			Size: api.ImageSize1x,
			Project: project{
				Desc: "Free Coupon",
				Type: api.Coupon,
			},
			NewKey: uuid.NewV4().String(),
		},
		{
			Name: "Coupon2x",
			Size: api.ImageSize2x,
			Project: project{
				Desc: "Free Coupon",
				Type: api.Coupon,
			},
			NewKey: uuid.NewV4().String(),
		},
		{
			Name: "Coupon3x",
			Size: api.ImageSize3x,
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

			conn, err := testConnection(ctx, t)
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

			p := api.NewProject(fakeString(), fakeString(), tc.Project.Desc, tc.Project.Type)

			err = db.SaveNewProject(ctx, u, p)
			if !assert.NoError(err) {
				return
			}

			err = db.SetBackgroundImage(ctx, tc.Size, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = db.SetFooterImage(ctx, tc.Size, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = db.SetIconImage(ctx, tc.Size, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = db.SetLogoImage(ctx, tc.Size, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			err = db.SetStripImage(ctx, tc.Size, tc.NewKey, p)
			if !assert.NoError(err) {
				return
			}

			loaded, err := db.LoadProject(ctx, u, p.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(p.ID, loaded.ID)

			switch tc.Size {
			case api.ImageSize3x:
				assert.Equal(tc.NewKey, loaded.BackgroundImage3x)
				assert.Equal(tc.NewKey, loaded.FooterImage3x)
				assert.Equal(tc.NewKey, loaded.IconImage3x)
				assert.Equal(tc.NewKey, loaded.LogoImage3x)
				assert.Equal(tc.NewKey, loaded.StripImage3x)
			case api.ImageSize2x:
				assert.Equal(tc.NewKey, loaded.BackgroundImage2x)
				assert.Equal(tc.NewKey, loaded.FooterImage2x)
				assert.Equal(tc.NewKey, loaded.IconImage2x)
				assert.Equal(tc.NewKey, loaded.LogoImage2x)
				assert.Equal(tc.NewKey, loaded.StripImage2x)
			default:
				assert.Equal(tc.NewKey, loaded.BackgroundImage)
				assert.Equal(tc.NewKey, loaded.FooterImage)
				assert.Equal(tc.NewKey, loaded.IconImage)
				assert.Equal(tc.NewKey, loaded.LogoImage)
				assert.Equal(tc.NewKey, loaded.StripImage)
			}
		})
	}
}

func TestProjectsPagination(t *testing.T) {
	testCases := []struct {
		Name     string
		Iters    int
		Projects int
		Limit    uint64
		HasNext  bool
		Expected int
	}{
		{
			Name:     "QueryFirstPage",
			Iters:    1,
			Projects: 20,
			Limit:    10,
			HasNext:  true,
			Expected: 10,
		},
		{
			Name:     "QueryFirstPageFull",
			Iters:    1,
			Projects: 7,
			Limit:    10,
			HasNext:  false,
			Expected: 7,
		},
		{
			Name:     "QuerySecondPage",
			Iters:    2,
			Projects: 20,
			Limit:    7,
			HasNext:  true,
			Expected: 7,
		},
		{
			Name:     "QuerySecondPageFull",
			Iters:    2,
			Projects: 20,
			Limit:    10,
			HasNext:  false,
			Expected: 10,
		},
		{
			Name:     "QueryAll",
			Iters:    3,
			Projects: 9,
			Limit:    3,
			HasNext:  false,
			Expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			conn, err := testConnection(ctx, t)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			user := api.NewUser(fakeUsername(), fakeEmail(), fakeString(), nil)
			err = db.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			for i := 0; i < tc.Projects; i++ {
				project := api.NewProject(fakeString(), fakeString(), fakeString(), api.Coupon)
				err = db.SaveNewProject(ctx, user, project)
				if !assert.NoError(err) {
					return
				}
			}

			var (
				opts        = api.NewPagingOptions(0, tc.Limit)
				output      *api.Projects
				hasNext     bool
				lastFetched int
			)

			for i := 0; i < tc.Iters; i++ {
				output, err = db.LoadProjects(ctx, user, opts)
				if !assert.NoError(err) {
					return
				}

				hasNext = output.Opts.HasNext()
				if hasNext {
					opts.Cursor = opts.Next
					opts.Next = 0
				}

				lastFetched = len(output.Data)
			}

			assert.Equal(tc.Expected, lastFetched)
			assert.Equal(tc.HasNext, hasNext)
		})
	}
}
