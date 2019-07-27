package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/dgrijalva/jwt-go"
)

const (
	// ServerIssuer is a JWT issuer name.
	ServerIssuer string = "OKPOCK"
	// ServerClaimsTTL is a default TTL for JWT token.
	ServerClaimsTTL time.Duration = 30 * 24 * time.Hour
)

var (
	serverSigningSecret                   = []byte("")
	serverSigningMethod jwt.SigningMethod = jwt.SigningMethodHS256
)

var (
	// ErrMissingStandardClaims returned when claims missing standard claims.
	ErrMissingStandardClaims = errors.New("claims: missing standard claims")
	// ErrMissingXSRFToken returned when claims missing csrf token.
	ErrMissingXSRFToken = errors.New("claims: missing csrf token")
	// ErrUnexpectedSigningMethod returned when claims has wrong signing method.
	ErrUnexpectedSigningMethod = errors.New("claims: unexpected signing method")
	// ErrInvalidClaims returned when token does not match claims.
	ErrInvalidClaims = errors.New("claims: token does not match claims")
)

// Claims ...
type Claims interface {
	jwt.Claims

	MarshalJWT() (string, error)
	UnmarshalJWT(string) error
}

// NewClaims creates a new user claims.
func NewClaims() *UserClaims {
	return &UserClaims{
		StandardClaims: &jwt.StandardClaims{
			Issuer:    ServerIssuer,
			ExpiresAt: jwt.TimeFunc().UTC().Add(ServerClaimsTTL).Unix(),
		},
	}
}

// UserClaims holds user claims.
type UserClaims struct {
	*jwt.StandardClaims

	CSRFToken string `json:"csrfToken"`
}

// WithUser updates claims' subject.
func (c *UserClaims) WithUser(u *api.User) *UserClaims {
	c.Subject = strconv.FormatInt(u.ID, 10)
	return c
}

// WithCSRFToken updates claims' XSRF token.
func (c *UserClaims) WithCSRFToken(token string) *UserClaims {
	c.CSRFToken = token
	return c
}

// Valid checks user claims.
func (c *UserClaims) Valid() error {
	if c.StandardClaims == nil {
		return ErrMissingStandardClaims
	}
	if c.CSRFToken == "" {
		return ErrMissingXSRFToken
	}
	return nil
}

// MarshalJWT generates JWT token.
func (c *UserClaims) MarshalJWT() (string, error) {
	if err := c.Valid(); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(serverSigningMethod, c)
	return token.SignedString(serverSigningSecret)
}

// UnmarshalJWT parses token string into user claims.
func (c *UserClaims) UnmarshalJWT(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, c, getSigningSecret)
	if err != nil {
		return err
	}
	if err := token.Claims.Valid(); err != nil {
		return err
	}
	tokenClaims, ok := token.Claims.(*UserClaims)
	if !ok {
		return ErrInvalidClaims
	}
	*c = *tokenClaims
	return nil
}

func getSigningSecret(token *jwt.Token) (interface{}, error) {
	if token.Method != serverSigningMethod {
		return nil, ErrUnexpectedSigningMethod
	}
	return serverSigningSecret, nil
}
