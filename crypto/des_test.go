package crypto

import (
	"testing"
)

func TestDESEncrypt(t *testing.T) {
	src := []byte("hello word")
	key := []byte("hello key")
	iv := []byte("@av3761L")
	encrypted, err := DESEncrypt(src, key, iv)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := DESDecrypt(encrypted, key, iv)
	if err != nil {
		t.Error(err)
	}
	if string(decrypted) != string(src) {
		t.Errorf("The decrypted result is: %s, which is not equal to the original string", string(decrypted))
	}
}
