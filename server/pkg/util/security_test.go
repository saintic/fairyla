package util

import (
	"strings"
	"testing"
)

func TestPasswd(t *testing.T) {
	pwd := "Abc123@$-*&"
	pwhash, err := GeneratePasswordHash(pwd)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !strings.HasPrefix(pwhash, "pbkdf2") {
		t.Fatal("hash password format error")
	}
	if !CheckPasswordHash(pwhash, pwd) {
		t.Fatal("CheckPasswordHash error (true)")
	}

	if CheckPasswordHash(pwhash, "ejaA5W4866") {
		t.Fatal("CheckPasswordHash error (false)")
	}
}
