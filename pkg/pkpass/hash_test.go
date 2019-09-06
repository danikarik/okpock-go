package pkpass_test

import (
	"io/ioutil"
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
