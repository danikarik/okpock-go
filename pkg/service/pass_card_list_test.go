package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestProjectPassCardsHandler(t *testing.T) {
	testCases := []struct {
		Name           string
		PassCardNumber int
	}{
		{Name: "EmptyPassCardList"},
		{Name: "NotEmptyPassCardList", PassCardNumber: 10},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			srv, err := initService(t)
			if !assert.NoError(err) {
				return
			}

			user := api.NewUser(fakeUsername(), fakeEmail(), fakePassword(), nil)
			err = srv.env.Auth.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			project := &api.Project{
				ID:               fakeID(),
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Description:      fakeString(),
				PassType:         api.Coupon,
			}

			err = srv.env.Logic.SaveNewProject(ctx, user, project)
			if !assert.NoError(err) {
				return
			}

			for i := 0; i < tc.PassCardNumber; i++ {
				passcard := fakePassCard(project)

				err := srv.env.Logic.SaveNewPassCard(ctx, project, passcard)
				if !assert.NoError(err) {
					return
				}
			}

			url := fmt.Sprintf("/projects/%d/cards", project.ID)
			req := authRequest(srv, user, newRequest("GET", url, nil, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				return
			}

			data := []*api.PassCardInfo{}
			err = unmarshalJSON(resp, &data)
			if !assert.NoError(err) {
				return
			}

			assert.Len(data, tc.PassCardNumber)
		})
	}
}

func TestProjectPassCardsWithBarcodeMessageHandler(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}

	user := api.NewUser(fakeUsername(), fakeEmail(), fakePassword(), nil)
	err = srv.env.Auth.SaveNewUser(ctx, user)
	if !assert.NoError(err) {
		return
	}

	project := &api.Project{
		ID:               fakeID(),
		Title:            fakeString(),
		OrganizationName: fakeString(),
		Description:      fakeString(),
		PassType:         api.Coupon,
	}

	err = srv.env.Logic.SaveNewProject(ctx, user, project)
	if !assert.NoError(err) {
		return
	}

	passcard := fakePassCard(project)
	passcard.Data.Barcodes = []*api.Barcode{
		&api.Barcode{
			Message:         fakeString(),
			Format:          "PKBarcodeFormatPDF417",
			MessageEncoding: "iso-8859-1",
		},
	}
	err = srv.env.Logic.SaveNewPassCard(ctx, project, passcard)
	if !assert.NoError(err) {
		return
	}

	url := fmt.Sprintf("/projects/%d/cards?barcode_message=%s", project.ID, passcard.Data.Barcodes[0].Message)
	req := authRequest(srv, user, newRequest("GET", url, nil, nil, nil))
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	if !assert.Equal(http.StatusOK, resp.StatusCode) {
		return
	}

	data := []*api.PassCardInfo{}
	err = unmarshalJSON(resp, &data)
	if !assert.NoError(err) {
		return
	}

	if assert.Len(data, 1) {
		loaded := data[0]
		assert.Equal(passcard.ID, loaded.ID)
		assert.Equal(passcard.Data, loaded.Data)
	}
}
