package secure_test

import (
	"testing"

	"github.com/danikarik/okpock/pkg/secure"
)

func TestSecureToken(t *testing.T) {
	num := 1000000
	generated := map[string]struct{}{}
	for i := 0; i < num; i++ {
		token := secure.Token()
		if _, ok := generated[token]; ok {
			t.Fatal("got duplicated token")
		}
	}
}
