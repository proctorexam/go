package procwise_test

import (
	"encoding/json"
	"io"
	"sort"
	"strings"
	"testing"
	"github.com/proctorexam/go/procwise"
)

func TestSignIn(t *testing.T) {
	c := procwise.NewClient(procwise.ClientOptions{
		Client: procwise.NewTestClient(),
	})
	_, err := c.SignIn("localhost", "user@proctorexam.com", "password")
	if err != nil {
		t.Error(err)
	}

	r := c.Client.(*procwise.TestClient).Requests[0]

	expectations := map[string]string{
		"POST": r.Method,
		"https://localhost.proctorexam.com:3001/users/sign_in.json": r.URL.String(),
		"application/json": r.Header.Get("Content-Type"),
		"email,password,password,user,user@proctorexam.com": func() string {
			b, _ := io.ReadAll(r.Body)
			var d map[string]map[string]string
			json.Unmarshal(b, &d)
			kv := make([]string, 0)
			for k, v := range d {
				kv = append(kv, k)
				for vk, vv := range v {
					kv = append(kv, vk, vv)
				}
			}
			sort.Strings(kv)
			return strings.Join(kv, ",")
		}(),
	}

	for expect, actual := range expectations {
		if expect != actual {
			t.Errorf("want %s got %s", expect, actual)
		}
	}
}

func TestFindForReview(t *testing.T) {
	c := procwise.NewClient(procwise.ClientOptions{
		Token:  "token",
		Secret: "secret",
		Client: procwise.NewTestClient(),
	})
	_, err := c.FindForReview("localhost", "cb28d066854085fe88002dece3e224a5")
	if err != nil {
		t.Error(err)
	}
	r := c.Client.(*procwise.TestClient).Requests[0]

	expectations := map[string]string{
		"GET": r.Method,
		"https://localhost.proctorexam.com:3001/api/v3/student_sessions/cb28d066854085fe88002dece3e224a5/reviewable": r.URL.String(),
		"Token token=token":           r.Header.Get("Authorization"),
		"application/vnd.procwise.v3": r.Header.Get("Accept"),
		"application/json":            r.Header.Get("Content-Type"),
		"id,nonce,signature,timestamp": func() string {
			b, _ := io.ReadAll(r.Body)
			var d map[string]any
			json.Unmarshal(b, &d)
			keys := make([]string, 0)
			for k := range d {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			return strings.Join(keys, ",")
		}(),
	}

	for expect, actual := range expectations {
		if expect != actual {
			t.Errorf("want %s got %s", expect, actual)
		}
	}
}
