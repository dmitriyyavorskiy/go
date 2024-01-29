package service

import (
	"regexp"
	"testing"
)

func TestSayHelloName(t *testing.T) {
	name := "Thor"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := SayHello("Thor")
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`SayHello("Thor") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestSayHelloEmpty(t *testing.T) {
	msg, err := SayHello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}
