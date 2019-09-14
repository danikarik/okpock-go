package memory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/secure"
	"github.com/danikarik/okpock/pkg/store/memory"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewPassCard(t *testing.T) {
	testCases := []struct {
		Name                 string
		PassType             api.PassType
		PassCard             *api.PassCard
		LoadBySerialNumber   bool
		LoadByBarcodeMessage bool
	}{
		{
			Name:     "CouponByID",
			PassType: api.Coupon,
			PassCard: &api.PassCard{
				FormatVersion: 1,
				SerialNumber:  fakeString(),
				TeamID:        fakeString(),
				Coupon: &api.PassStructure{
					AuxiliaryFields: []*api.Field{
						&api.Field{
							Key:        "expires",
							Label:      "EXPIRES",
							Value:      "2020-04-24T10:00-05:00",
							IsRelative: true,
							DateStyle:  api.PKDateStyleShort,
						},
					},
					BackFields: []*api.Field{
						&api.Field{
							Key:   "offer",
							Label: "Any premium dog food",
							Value: "20% off",
						},
					},
				},
				Barcodes: []*api.Barcode{
					&api.Barcode{
						Message:         fakeString(),
						Format:          api.PKBarcodeFormatPDF417,
						MessageEncoding: "iso-8859-1",
					},
				},
				AuthenticationToken: secure.Token(),
				WebServiceURL:       "https://okpock.com",
			},
		},
		{
			Name:     "CouponBySerialNumber",
			PassType: api.Coupon,
			PassCard: &api.PassCard{
				FormatVersion: 1,
				SerialNumber:  fakeString(),
				TeamID:        fakeString(),
				Coupon: &api.PassStructure{
					AuxiliaryFields: []*api.Field{
						&api.Field{
							Key:        "expires",
							Label:      "EXPIRES",
							Value:      "2020-04-24T10:00-05:00",
							IsRelative: true,
							DateStyle:  api.PKDateStyleShort,
						},
					},
					BackFields: []*api.Field{
						&api.Field{
							Key:   "offer",
							Label: "Any premium dog food",
							Value: "20% off",
						},
					},
				},
				Barcodes: []*api.Barcode{
					&api.Barcode{
						Message:         fakeString(),
						Format:          api.PKBarcodeFormatPDF417,
						MessageEncoding: "iso-8859-1",
					},
				},
				AuthenticationToken: secure.Token(),
				WebServiceURL:       "https://okpock.com",
			},
			LoadBySerialNumber: true,
		},
		{
			Name:     "CouponByBarcodeMessage",
			PassType: api.Coupon,
			PassCard: &api.PassCard{
				FormatVersion: 1,
				SerialNumber:  fakeString(),
				TeamID:        fakeString(),
				Coupon: &api.PassStructure{
					AuxiliaryFields: []*api.Field{
						&api.Field{
							Key:        "expires",
							Label:      "EXPIRES",
							Value:      "2020-04-24T10:00-05:00",
							IsRelative: true,
							DateStyle:  api.PKDateStyleShort,
						},
					},
					BackFields: []*api.Field{
						&api.Field{
							Key:   "offer",
							Label: "Any premium dog food",
							Value: "20% off",
						},
					},
				},
				Barcodes: []*api.Barcode{
					&api.Barcode{
						Message:         fakeString(),
						Format:          api.PKBarcodeFormatPDF417,
						MessageEncoding: "iso-8859-1",
					},
				},
				AuthenticationToken: secure.Token(),
				WebServiceURL:       "https://okpock.com",
			},
			LoadByBarcodeMessage: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			db := memory.New()

			assert := assert.New(t)

			user := api.NewUser(fakeUsername(), fakeString(), fakeString(), nil)
			err := db.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			project := &api.Project{
				ID:               fakeID(),
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Description:      fakeString(),
				PassType:         tc.PassType,
			}
			err = db.SaveNewProject(ctx, user, project)
			if !assert.NoError(err) {
				return
			}

			tc.PassCard.Description = project.Description
			tc.PassCard.OrganizationName = project.OrganizationName
			tc.PassCard.PassTypeID = fmt.Sprintf("pass.okpock.com.%s", project.PassType)

			passcard := api.NewPassCardInfo(tc.PassCard)
			err = db.SaveNewPassCard(ctx, project, passcard)
			if !assert.NoError(err) {
				return
			}

			var loaded *api.PassCardInfo
			{
				if tc.LoadBySerialNumber {
					loaded, err = db.LoadPassCardBySerialNumber(ctx, project, passcard.Data.SerialNumber)
				} else {
					loaded, err = db.LoadPassCard(ctx, project, passcard.ID)
				}
			}
			if !assert.NoError(err) {
				return
			}
			if !assert.NotNil(loaded) {
				return
			}

			if !assert.Equal(passcard, loaded) {
				return
			}

			var loadedPassCards []*api.PassCardInfo
			{
				if tc.LoadByBarcodeMessage {
					loadedPassCards, err = db.LoadPassCardsByBarcodeMessage(ctx, project, passcard.Data.Barcodes[0].Message)
				} else {
					loadedPassCards, err = db.LoadPassCards(ctx, project)
				}
			}
			if !assert.NoError(err) {
				return
			}
			if !assert.NotNil(loadedPassCards) {
				return
			}

			if assert.Len(loadedPassCards, 1) {
				loaded = loadedPassCards[0]
				assert.Equal(passcard, loaded)
			}
		})
	}
}

func TestUpdatePassCard(t *testing.T) {
	testCases := []struct {
		Name      string
		PassType  api.PassType
		PassCard  *api.PassCard
		Locations []*api.Location
	}{
		{
			Name:     "Coupon",
			PassType: api.Coupon,
			PassCard: &api.PassCard{
				FormatVersion: 1,
				SerialNumber:  fakeString(),
				TeamID:        fakeString(),
				Coupon: &api.PassStructure{
					AuxiliaryFields: []*api.Field{
						&api.Field{
							Key:        "expires",
							Label:      "EXPIRES",
							Value:      "2020-04-24T10:00-05:00",
							IsRelative: true,
							DateStyle:  api.PKDateStyleShort,
						},
					},
					BackFields: []*api.Field{
						&api.Field{
							Key:   "offer",
							Label: "Any premium dog food",
							Value: "20% off",
						},
					},
				},
				Barcodes: []*api.Barcode{
					&api.Barcode{
						Message:         fakeString(),
						Format:          api.PKBarcodeFormatPDF417,
						MessageEncoding: "iso-8859-1",
					},
				},
				AuthenticationToken: secure.Token(),
				WebServiceURL:       "https://okpock.com",
			},
			Locations: []*api.Location{
				&api.Location{
					Longitude: -122.3748889,
					Latitude:  37.6189722,
				},
				&api.Location{
					Longitude: -122.03118,
					Latitude:  37.33182,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			db := memory.New()

			assert := assert.New(t)

			user := api.NewUser(fakeUsername(), fakeString(), fakeString(), nil)
			err := db.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			project := &api.Project{
				ID:               fakeID(),
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Description:      fakeString(),
				PassType:         tc.PassType,
			}
			err = db.SaveNewProject(ctx, user, project)
			if !assert.NoError(err) {
				return
			}

			tc.PassCard.Description = project.Description
			tc.PassCard.OrganizationName = project.OrganizationName
			tc.PassCard.PassTypeID = fmt.Sprintf("pass.okpock.com.%s", project.PassType)

			passcard := api.NewPassCardInfo(tc.PassCard)
			err = db.SaveNewPassCard(ctx, project, passcard)
			if !assert.NoError(err) {
				return
			}

			newData := passcard.Data
			newData.Locations = tc.Locations
			err = db.UpdatePassCard(ctx, newData, passcard)
			if !assert.NoError(err) {
				return
			}

			loaded, err := db.LoadPassCard(ctx, project, passcard.ID)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(passcard, loaded)
			assert.Equal(loaded.Data.Locations, tc.Locations)
		})
	}
}
