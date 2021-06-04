package util

import (
	"strings"
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

	if IsImage("a.jpg") != true || IsImage("a.xx") == true {
		t.Fatal("IsImage error")
	}
	if IsVideo("a.mp4") != true || IsVideo("a.xx") == true {
		t.Fatal("IsVideo error")
	}

	st := []string{"a", "b", "c"}
	st1 := DeleteSlice(st, "a")
	if strings.Join(st1, "") != "bc" {
		t.Fatal("Delete slice error(1)")
	}
	st2 := DeleteSlice(st, "b")
	if strings.Join(st2, "") != "ac" {
		t.Fatal("Delete slice error(2)")
	}
	st3 := DeleteSlice(st, "c")
	if strings.Join(st3, "") != "ab" {
		t.Fatal("Delete slice error(3)")
	}
	st4 := DeleteSlice(st, "d")
	if strings.Join(st4, "") != "abc" {
		t.Fatal("Delete slice error(4)")
	}

	okMails := []string{"a@b.com", "x-y@z.io", "hello.world@linux.org"}
	errMails := []string{"a@b.", "xx.com", "123456"}
	for _, v := range okMails {
		if IsEmail(v) != true {
			t.Fatalf("IsEmail should be true: %s\n", v)
		}
	}
	for _, v := range errMails {
		if IsEmail(v) == true {
			t.Fatalf("IsEmail should be false: %s\n", v)
		}
	}
}
