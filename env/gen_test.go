package env_test

import (
	"testing"

	"github.com/proctorexam/go/env"
)

type tWriter struct {
	b string
}

func (w *tWriter) Write(b []byte) (int, error) {
	w.b += string(b)
	return len(b), nil
}

func TestGen(t *testing.T) {
	w := &tWriter{}
	env.Gen(w, "env", []env.Var{
		{Name: "GET", Method: "Get"},
		{Name: "FETCH", Method: "Fetch", Default: "fetch"},
		{Name: "MUST", Method: "Must"},
	})

	if genExpect != w.b {
		t.Errorf("gen output does not match: want %s, got %s", genExpect, w.b)
	}
}

var genExpect = `package env

import e "github.com/proctorexam/go/env"

var GET string
var FETCH string
var MUST string

func init() {
	GET = e.Get("GET")
	FETCH = e.Fetch("FETCH","fetch")
	MUST = e.Must("MUST")
}
`
