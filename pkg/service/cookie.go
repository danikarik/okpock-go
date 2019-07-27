package service

import (
	"errors"
	"net/http"
	"time"
)

const (
	// CookieDomain domain name used for cookies.
	CookieDomain string = ".okpock.com"
	// TokenCookieName used for JWT token.
	TokenCookieName string = "okpocktok"
)

// ErrInvalidToken returns when cookie is missing or malformed.
var ErrInvalidToken = errors.New("cookie: missing or malformed access token")

func getClaims(r *http.Request, c Claims) error {
	cookie, err := r.Cookie(TokenCookieName)
	if err == http.ErrNoCookie {
		return ErrInvalidToken
	} else if err != nil {
		return err
	}
	tokenString := cookie.Value
	if tokenString == "" {
		return ErrInvalidToken
	}
	return c.UnmarshalJWT(tokenString)
}

func setClaimsCookie(w http.ResponseWriter, c Claims) error {
	tokenString, err := c.MarshalJWT()
	if err != nil {
		return err
	}
	http.SetCookie(w, tokenCookie(tokenString))
	return nil
}

func tokenCookie(tokenString string) *http.Cookie {
	return &http.Cookie{
		Name:     TokenCookieName,
		Domain:   CookieDomain,
		Path:     "/",
		Expires:  time.Now().UTC().Add(ServerClaimsTTL),
		Secure:   true,
		HttpOnly: true,
		Value:    tokenString,
	}
}

func clearCookies(w http.ResponseWriter) error {
	cookie := tokenCookie("")
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
	cookie.Value = ""
	http.SetCookie(w, cookie)
	return nil
}
