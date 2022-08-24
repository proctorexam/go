package greeter_test

import (
	"testing"

	"github.com/proctorexam/go/greeter"
)

func TestHello(t *testing.T) {
	if greeter.Hello() != "hello" {
		t.Error("hello is not hello")
	}
}
