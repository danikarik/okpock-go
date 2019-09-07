package pkpass_test

import (
	"io/ioutil"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/pkpass"
	"github.com/stretchr/testify/assert"
)

func TestPKCS7Sign(t *testing.T) {
	env, err := env.NewLookup("COUPON_P12")
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

			signer := pkpass.NewSigner(env.Get("COUPON_P12"))

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
