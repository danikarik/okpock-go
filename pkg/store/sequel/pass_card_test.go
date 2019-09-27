package sequel_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/secure"
	"github.com/danikarik/okpock/pkg/store/sequel"
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
			assert := assert.New(t)

			conn, err := testConnection(ctx, t)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			user := api.NewUser(fakeUsername(), fakeString(), fakeString(), nil)
			err = db.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			project := &api.Project{
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Description:      fakeString(),
				PassType:         tc.PassType,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
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

			assert.Equal(passcard.ID, loaded.ID)
			assert.Equal(passcard.Data, loaded.Data)

			var loadedPassCards *api.PassCardInfoList
			{
				if tc.LoadByBarcodeMessage {
					loadedPassCards, err = db.LoadPassCardsByBarcodeMessage(ctx, project, passcard.Data.Barcodes[0].Message, nil)
				} else {
					loadedPassCards, err = db.LoadPassCards(ctx, project, nil)
				}
			}
			if !assert.NoError(err) {
				return
			}
			if !assert.NotNil(loadedPassCards) {
				return
			}

			if assert.Len(loadedPassCards.Data, 1) {
				loaded = loadedPassCards.Data[0]
				assert.Equal(passcard.ID, loaded.ID)
				assert.Equal(passcard.Data, loaded.Data)
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
			assert := assert.New(t)

			conn, err := testConnection(ctx, t)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			user := api.NewUser(fakeUsername(), fakeString(), fakeString(), nil)
			err = db.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			project := &api.Project{
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Description:      fakeString(),
				PassType:         tc.PassType,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
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

			assert.Equal(passcard.ID, loaded.ID)
			assert.Equal(passcard.Data, loaded.Data)
			assert.Equal(loaded.Data.Locations, tc.Locations)
		})
	}
}

func TestPassCardsPagination(t *testing.T) {
	testCases := []struct {
		Name              string
		UseBarcodeMessage bool
		Iters             int
		PassCards         int
		Limit             uint64
		HasNext           bool
		Expected          int
	}{
		{
			Name:      "QueryFirstPage",
			Iters:     1,
			PassCards: 20,
			Limit:     10,
			HasNext:   true,
			Expected:  10,
		},
		{
			Name:              "QueryFirstPageWithBarcode",
			UseBarcodeMessage: true,
			Iters:             1,
			PassCards:         20,
			Limit:             10,
			HasNext:           true,
			Expected:          10,
		},
		{
			Name:      "QueryFirstPageFull",
			Iters:     1,
			PassCards: 7,
			Limit:     10,
			HasNext:   false,
			Expected:  7,
		},
		{
			Name:              "QueryFirstPageFullWithBarcode",
			UseBarcodeMessage: true,
			Iters:             1,
			PassCards:         7,
			Limit:             10,
			HasNext:           false,
			Expected:          7,
		},
		{
			Name:      "QuerySecondPage",
			Iters:     2,
			PassCards: 20,
			Limit:     7,
			HasNext:   true,
			Expected:  7,
		},
		{
			Name:              "QuerySecondPageWithBarcode",
			UseBarcodeMessage: true,
			Iters:             2,
			PassCards:         20,
			Limit:             7,
			HasNext:           true,
			Expected:          7,
		},
		{
			Name:      "QuerySecondPageFull",
			Iters:     2,
			PassCards: 20,
			Limit:     10,
			HasNext:   false,
			Expected:  10,
		},
		{
			Name:              "QuerySecondPageFullWithBarcode",
			UseBarcodeMessage: true,
			Iters:             2,
			PassCards:         20,
			Limit:             10,
			HasNext:           false,
			Expected:          10,
		},
		{
			Name:      "QueryAll",
			Iters:     3,
			PassCards: 9,
			Limit:     3,
			HasNext:   false,
			Expected:  3,
		},
		{
			Name:              "QueryAllWithBarcode",
			UseBarcodeMessage: true,
			Iters:             3,
			PassCards:         9,
			Limit:             3,
			HasNext:           false,
			Expected:          3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			member := fakeString()

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

			project := api.NewProject(fakeString(), fakeString(), fakeString(), api.Coupon)
			err = db.SaveNewProject(ctx, user, project)
			if !assert.NoError(err) {
				return
			}

			for i := 0; i < tc.PassCards; i++ {
				data := &api.PassCard{
					Description:      project.Description,
					FormatVersion:    1,
					OrganizationName: project.OrganizationName,
					PassTypeID:       fmt.Sprintf("pass.okpock.com.%s", project.PassType),
					SerialNumber:     fakeString(),
					TeamID:           fakeString(),
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
							Message:         member,
							Format:          api.PKBarcodeFormatPDF417,
							MessageEncoding: "iso-8859-1",
						},
					},
					AuthenticationToken: secure.Token(),
					WebServiceURL:       "https://okpock.com",
				}

				passcard := api.NewPassCardInfo(data)
				err = db.SaveNewPassCard(ctx, project, passcard)
				if !assert.NoError(err) {
					return
				}
			}

			var (
				opts        = api.NewPagingOptions(0, tc.Limit)
				output      *api.PassCardInfoList
				hasNext     bool
				lastFetched int
			)

			for i := 0; i < tc.Iters; i++ {
				if tc.UseBarcodeMessage {
					output, err = db.LoadPassCardsByBarcodeMessage(ctx, project, member, opts)
				} else {
					output, err = db.LoadPassCards(ctx, project, opts)
				}
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
