package pkpass_test

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/pkpass"
	"github.com/stretchr/testify/assert"
)

func skipTest(t *testing.T) {
	if v, ok := os.LookupEnv("SKIP_SIGNER_TEST"); ok {
		skip, err := strconv.ParseBool(v)
		if err == nil && skip {
			t.Skip(`skip test: SKIP_SIGNER_TEST is present`)
		}
	}
}

func TestPKCS7Sign(t *testing.T) {
	skipTest(t)

	env, err := env.NewLookup(
		"TEST_CERTIFICATES_ROOT_CERT",
		"TEST_CERTIFICATES_COUPON_PATH",
		"TEST_CERTIFICATES_COUPON_PASS",
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

			root, err := ioutil.ReadFile(env.Get("TEST_CERTIFICATES_ROOT_CERT"))
			if !assert.NoError(err) {
				return
			}

			signing, err := ioutil.ReadFile(env.Get("TEST_CERTIFICATES_COUPON_PATH"))
			if !assert.NoError(err) {
				return
			}

			signer, err := pkpass.NewSigner(root, signing, env.Get("TEST_CERTIFICATES_COUPON_PASS"))
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
