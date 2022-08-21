package env

import (
	"bytes"
	"fmt"
	"os"
)

var newLine = []byte("\n")
var eqSign = []byte("=")

// Must lookups for an environment variable for given name,
// panics if it's not defined, returns it's value otherwise
func Must(n string) string {
	v, ok := os.LookupEnv(n)
	if !ok {
		panic(fmt.Errorf("missing env var: %s", n))
	}
	return v
}

// Get returns the value of environment variable for given name
// If variable is not set, return value is empty string
func Get(n string) string {
	return os.Getenv(n)
}

// Fetch is same with Get, but accepts default value
func Fetch(n string, d ...string) string {
	v, ok := os.LookupEnv(n)
	if !ok && len(d) > 0 {
		return d[0]
	}
	return v
}

// Load sets variables from a file for given filepath
func Load(f string) error {
	b, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	for _, li := range bytes.Split(b, newLine) {
		li = bytes.TrimSpace(li)
		kv := bytes.SplitN(li, eqSign, 2)
		if len(kv) != 2 {
			continue
		}
		os.Setenv(string(kv[0]), string(kv[1]))
	}
	return nil
}
