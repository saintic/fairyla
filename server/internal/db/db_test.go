package db

import (
	"os"
	"testing"

	"tcw.im/ufc"
)

var c *Conn

func raise(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func getConn(t *testing.T) {
	if c == nil {
		rawurl := os.Getenv("redis_url")
		if rawurl == "" {
			t.SkipNow()
		} else {
			db, err := New(rawurl)
			raise(t, err)
			c = db
		}
	}
}

func TestDBString(t *testing.T) {
	getConn(t)

	k := "test"

	ok, err := c.Set(k, "value")
	raise(t, err)
	if !ok {
		t.Fatal("set error")
	}

	v, err := c.Get(k)
	raise(t, err)
	if v != "value" {
		t.Fatal("get error")
	}

	ks, err := c.Keys("*")
	raise(t, err)
	if !ufc.StrInSlice(c.Prefix+k, ks) {
		t.Fatal("keys error")
	}

	typ, err := c.Type(k)
	raise(t, err)
	if typ != "string" {
		t.Fatal("string type error")
	}

	has, err := c.Exsits(k)
	raise(t, err)
	if !has {
		t.Fatal("exists error")
	}

	_, err = c.Del(k)
	raise(t, err)
	has, _ = c.Exsits(k)
	if has {
		t.Fatal("del error")
	}
}
