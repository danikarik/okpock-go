package service

import (
	"os"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestUserClaims(t *testing.T) {
	if v, ok := os.LookupEnv("TEST_SERVER_SECRET"); !ok {
		t.Skip(`"TEST_SERVER_SECRET" is not present`)
	} else {
		serverSigningSecret = []byte(v)
	}

	assert := assert.New(t)

	c := NewClaims()
	err := c.Valid()
	assert.Error(err)

	u, err := api.NewUser("test", "test@example.com", "test", nil)
	assert.NoError(err)
	c = c.WithUser(u).WithCSRFToken("token")
	err = c.Valid()
	assert.NoError(err)

	tokenString, err := c.MarshalJWT()
	assert.NoError(err)
	assert.NotEmpty(tokenString)

	c2 := NewClaims()
	err = c2.UnmarshalJWT(tokenString)
	assert.NoError(err)
	assert.Equal(c.Subject, c2.Subject)
	assert.Equal(c.CSRFToken, c2.CSRFToken)
}
