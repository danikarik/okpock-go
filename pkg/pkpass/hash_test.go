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
			Name:     "TextFile",
			Path:     "testdata/test.txt",
			Expected: "740c172dc6bd2f9262eb3e19d080b7f106249899",
		},
		{
			Name:     "ImageFile",
			Path:     "testdata/gopher.jpg",
			Expected: "4e8069b789897df23449b8a1bbc812232e36a7d3",
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
		Name  string
		Paths []string
	}{
		{
			Name: "SimpleOne",
			Paths: []string{
				"testdata/gopher.jpg",
				"testdata/test.txt",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			hashes := make(map[string]string)
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

				hash, err := pkpass.HashFile(data)
				if !assert.NoError(err) {
					return
				}

				hashes[fi.Name()] = hash

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

			for filename, hash := range loadedManifest {
				v, ok := hashes[filename]
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
