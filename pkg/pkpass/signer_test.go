package pkpass_test

import (
	"io/ioutil"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/pkpass"
	"github.com/stretchr/testify/assert"
)

func TestPKCS7Sign(t *testing.T) {
	env, err := env.NewLookup(
		"WWDR_CERTIFICATE",
		"COUPON_CERTIFICATE",
		"COUPON_PASSWORD",
	)
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

			root, err := ioutil.ReadFile(env.Get("WWDR_CERTIFICATE"))
			if !assert.NoError(err) {
				return
			}

			signing, err := ioutil.ReadFile(env.Get("COUPON_CERTIFICATE"))
			if !assert.NoError(err) {
				return
			}

			signer, err := pkpass.NewSigner(root, signing, env.Get("COUPON_PASSWORD"))
			if !assert.NoError(err) {
				return
			}

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
