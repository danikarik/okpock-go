package pkpass_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/danikarik/okpock/pkg/pkpass"
	"github.com/stretchr/testify/assert"
)

func TestZip(t *testing.T) {
	testCases := []struct {
		Name  string
		Paths []string
	}{
		{
			Name: "Coupon",
			Paths: []string{
				"testdata/coupon.pass/icon.png",
				"testdata/coupon.pass/icon@2x.png",
				"testdata/coupon.pass/pass.json",
				"testdata/coupon.pass/logo.png",
				"testdata/coupon.pass/logo@2x.png",
				"testdata/coupon.pass/manifest.json",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

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

			zipContent, err := pkpass.Zip(files...)
			if !assert.NoError(err) {
				return
			}

			unzippedFiles, err := pkpass.Unzip(zipContent)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(unzippedFiles, files)
		})
	}
}
