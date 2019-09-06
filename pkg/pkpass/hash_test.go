package pkpass_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/danikarik/okpock/pkg/pkpass"
	"github.com/stretchr/testify/assert"
)

func TestHashFile(t *testing.T) {
	testCases := []struct {
		Name     string
		Path     string
		Expected string
	}{
		{
			Name:     "CouponIcon",
			Path:     "testdata/coupon.pass/icon.png",
			Expected: "f8a2bb1b52c426275312c98c626d5be92758170e",
		},
		{
			Name:     "CouponIcon2x",
			Path:     "testdata/coupon.pass/icon@2x.png",
			Expected: "4204eafa4ac2df2339cf3308a2b0ecd228732589",
		},
		{
			Name:     "CouponPassJson",
			Path:     "testdata/coupon.pass/pass.json",
			Expected: "ac5eccd991c295c58d7daf0675c05f67973a6321",
		},
		{
			Name:     "CouponLogo",
			Path:     "testdata/coupon.pass/logo.png",
			Expected: "2147d5e2561b98bc2d9653c51349947d6ddc419d",
		},
		{
			Name:     "CouponLogo2x",
			Path:     "testdata/coupon.pass/logo@2x.png",
			Expected: "b98b0504f4f067de4f7a6c1e95df8a78024dc3bb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			content, err := ioutil.ReadFile(tc.Path)
			if !assert.NoError(err) {
				return
			}

			hash, err := pkpass.HashFile(content)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(tc.Expected, hash)
		})
	}
}

func TestManifest(t *testing.T) {
	testCases := []struct {
		Name     string
		Paths    []string
		Expected string
	}{
		{
			Name: "Coupon",
			Paths: []string{
				"testdata/coupon.pass/icon.png",
				"testdata/coupon.pass/icon@2x.png",
				"testdata/coupon.pass/pass.json",
				"testdata/coupon.pass/logo.png",
				"testdata/coupon.pass/logo@2x.png",
			},
			Expected: "testdata/coupon.pass/manifest.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			manifestData, err := ioutil.ReadFile(tc.Expected)
			if !assert.NoError(err) {
				return
			}

			files := make([]pkpass.File, len(tc.Paths))
			for i, path := range tc.Paths {
				file, err := os.Open(path)
				if !assert.NoError(err) {
					return
				}
				defer file.Close()

				fi, err := file.Stat()
				if !assert.NoError(err) {
					return
				}

				data, err := ioutil.ReadAll(file)
				if !assert.NoError(err) {
					return
				}

				files[i] = pkpass.File{
					Name: fi.Name(),
					Data: data,
				}
			}

			manifest, err := pkpass.CreateManifest(files...)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(pkpass.ManifestFilename, manifest.Name)

			loadedManifest := pkpass.Manifest{}
			err = json.Unmarshal(manifest.Data, &loadedManifest)
			if !assert.NoError(err) {
				return
			}

			expectedManifest := pkpass.Manifest{}
			err = json.Unmarshal(manifestData, &expectedManifest)
			if !assert.NoError(err) {
				return
			}

			for filename, hash := range loadedManifest {
				v, ok := expectedManifest[filename]
				if !assert.True(ok) {
					return
				}
				if !assert.Equal(v, hash) {
					return
				}
			}
		})
	}
}
