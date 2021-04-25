package util

import (
	"testing"
)

func TestUtil(t *testing.T) {
	okNames := []string{"test", "a1", "x-", "y_", "z-hello", "b0123", "adsatgr"}
	errNames := []string{"a", "-", "-a", "aB", "_b", "0a", "abc@", "a#", "Test"}
	for _, v := range okNames {
		if IsName(v) != true {
			t.Fatalf("IsName should be true: %s\n", v)
		}
	}
	for _, v := range errNames {
		if IsName(v) == true {
			t.Fatalf("IsName should be false: %s\n", v)
		}
	}

	okURLs := []string{
		"https://abc.com", "http://127.0.0.1", "https://x.com:8443",
		"https://user:pass@github.com/user/repo", "https://-x.y.z",
	}
	errURLs := []string{
		"local", "127.0.0.1", "http:///abc", "x//y.z", "x://y.z", "a.com",
		"ftp://ftp.hello.world", "redis://localhost",
	}
	for _, v := range okURLs {
		if IsValidURL(v) != true {
			t.Fatalf("IsValidURL should be true: %s\n", v)
		}
	}
	for _, v := range errURLs {
		if IsValidURL(v) == true {
			t.Fatalf("IsValidURL should be false: %s\n", v)
		}
	}
}
