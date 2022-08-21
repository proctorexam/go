package env_test

import (
	"os"
	"testing"

	"github.com/proctorexam/go/env"
)

func TestMust(t *testing.T) {
	os.Setenv("must", "ok")
	if env.Must("must") != "ok" {
		t.Error("env.Must() not returned ok")
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Error("Must not paniced")
		}
	}()
	env.Must("nope")
}

func TestGet(t *testing.T) {
	os.Setenv("test", "ok")
	if env.Get("test") != "ok" {
		t.Error("env.Get() not returned ok")
	}
}

func TestFetch(t *testing.T) {
	if env.Fetch("fetch", "ok") != "ok" {
		t.Error("env.Fetch() not returned ok")
	}
	os.Setenv("fetch", "ok")
	if env.Fetch("fetch") != "ok" {
		t.Error("env.Fetch() not returned ok")
	}
}

func TestLoad(t *testing.T) {
	err := env.Load("testdata/.env")
	if err != nil {
		t.Fatal(err)
	}
	if env.Get("EMPTY") != "" {
		t.Error("env.Load() did not set EMPTY")
	}
	if env.Get("SIMPLE") != "simple" {
		t.Error("env.Load() did not set SIMPLE")
	}
	if env.Get("URL") != "https://user@password:domain.com:8000?q=$1#hash" {
		t.Error("env.Load() did not set URL")
	}
}
