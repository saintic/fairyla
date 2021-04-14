package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestConf(t *testing.T) {
	data := []byte(`
    gn = global
    [project]
    latest = master
    `)
	f := filepath.Join(os.TempDir(), "_fairyla_conf_test.ini")
	err := ioutil.WriteFile(f, data, 0644)
	if err != nil {
		t.Fatal("write test file error")
	}

	cfg, err := New(f)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.GetKey("gn") != "global" {
		t.Fatal("get key error")
	}

	if cfg.MustKey("non", "dft") != "dft" {
		t.Fatal("must get key error")
	}

	if cfg.GetSecKey("project", "latest") != "master" {
		t.Fatal("get section key error")
	}
	if cfg.MustSecKey("nonsec", "nonkey", "dft") != "dft" {
		t.Fatal("must get section key error")
	}
}
