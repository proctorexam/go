package env_test

import (
	"github.com/proctorexam/go/procwise/env"
	"testing"
)

func TestVars(t *testing.T) {
	cases := map[string]string{
		"PE_USER":          env.PE_USER,
		"PE_PASSWORD":      env.PE_PASSWORD,
	}
	for k, v := range cases {
		if v == "" {
			t.Errorf("env var %s is empty", k)
		}
	}
}
