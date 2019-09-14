package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestDictionaryHandlers(t *testing.T) {
	testCases := []struct {
		Name     string
		Route    string
		Expected []string
	}{
		{
			Name:  "PassTypes",
			Route: "passtypes",
			Expected: []string{
				string(api.BoardingPass),
				string(api.Coupon),
				string(api.EventTicket),
				string(api.Generic),
				string(api.StoreCard),
			},
		},
		{
			Name:  "DetectorTypes",
			Route: "detectortypes",
			Expected: []string{
				api.PKDataDetectorTypePhoneNumber,
				api.PKDataDetectorTypeLink,
				api.PKDataDetectorTypeAddress,
				api.PKDataDetectorTypeCalendarEvent,
			},
		},
		{
			Name:  "TextAlignment",
			Route: "textalignment",
			Expected: []string{
				api.PKTextAlignmentLeft,
				api.PKTextAlignmentCenter,
				api.PKTextAlignmentRight,
				api.PKTextAlignmentNatural,
			},
		},
		{
			Name:  "DateStyle",
			Route: "datestyle",
			Expected: []string{
				api.PKDateStyleNone,
				api.PKDateStyleShort,
				api.PKDateStyleMedium,
				api.PKDateStyleLong,
				api.PKDateStyleFull,
			},
		},
		{
			Name:  "NumberStyle",
			Route: "numberstyle",
			Expected: []string{
				api.PKNumberStyleDecimal,
				api.PKNumberStylePercent,
				api.PKNumberStyleScientific,
				api.PKNumberStyleSpellOut,
			},
		},
		{
			Name:  "TransitType",
			Route: "transittype",
			Expected: []string{
				api.PKTransitTypeAir,
				api.PKTransitTypeBoat,
				api.PKTransitTypeBus,
				api.PKTransitTypeGeneric,
				api.PKTransitTypeTrain,
			},
		},
		{
			Name:  "BarcodeFormat",
			Route: "barcodeformat",
			Expected: []string{
				api.PKBarcodeFormatQR,
				api.PKBarcodeFormatPDF417,
				api.PKBarcodeFormatAztec,
				api.PKBarcodeFormatCode128,
			},
		},
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

			req := authRequest(srv, user, newRequest("GET", "/dictionary/"+tc.Route, nil, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				return
			}

			var data = M{}
			err = unmarshalJSON(resp, &data)
			if !assert.NoError(err) {
				return
			}

			for _, item := range data["data"].([]interface{}) {
				assert.Contains(tc.Expected, item)
			}
		})
	}
}
