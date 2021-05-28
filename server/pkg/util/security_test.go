package util

import (
	"strings"
	"testing"
)

func TestPasswd(t *testing.T) {
	pwd := "Abc123@$-*&"
	pwhash := GeneratePasswordHash(pwd)
	if !strings.HasPrefix(pwhash, "pbkdf2") {
		t.Fatal("hash password format error")
	}
	if !CheckPasswordHash(pwhash, pwd) {
		t.Fatal("CheckPasswordHash error (true)")
	}

	if CheckPasswordHash(pwhash, "ejaA5W4866") {
		t.Fatal("CheckPasswordHash error (false)")
	}

	if !CheckPasswordHash("pbkdf2:sha256:150000$YvakFsKw$98daf134870a2578202f09fb3b9f5564bad38541fa66cbd4c296962a5ab836b4", "123456") {
		t.Fatal("CheckPasswordHash error (old)")
	}

	// test werkzeug generate_password_hash
	pypwhash := "pbkdf2:sha256:150000$XqnEnB37$cf2bffd0ccf0ae5d5c7c9486a6c3dc8f4add525f6faf9c4d85c91e98bb360b86"
	if !CheckPasswordHash(pypwhash, "123456") {
		t.Fatal("CheckPasswordHash error (werkzeug)")
	}
}
