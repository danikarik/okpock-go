package pkpass_test

import (
	"io/ioutil"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/pkpass"
	"github.com/stretchr/testify/assert"
)

var requiredVars = []string{
	"WWDR_CERTIFICATE",
	"COUPON_CERTIFICATE",
	"COUPON_KEY",
	"COUPON_PASSWORD",
}

func TestOpenSSLSign(t *testing.T) {
	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	testCases := []struct {
		Name     string
		Path     string
		Expected string
	}{
		{
			Name:     "Coupon",
			Path:     "testdata/coupon.pass/manifest.json",
			Expected: "testdata/coupon.pass/signature",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			signer := pkpass.NewOpenSSL(
				env.Get("WWDR_CERTIFICATE"),
				env.Get("COUPON_CERTIFICATE"),
				env.Get("COUPON_KEY"),
				env.Get("COUPON_PASSWORD"),
			)

			data, err := ioutil.ReadFile(tc.Path)
			if !assert.NoError(err) {
				return
			}

			signatureData, err := ioutil.ReadFile(tc.Expected)
			if !assert.NoError(err) {
				return
			}

			signatureFile, err := signer.Sign(data)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(pkpass.SignatureFilename, signatureFile.Name)
			assert.Len(signatureFile.Data, len(signatureData))
		})
	}
}
