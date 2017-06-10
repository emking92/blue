package utils

import (
	"testing"
)

func TestStrSubstitutor(t *testing.T) {
	ss := make(StrSubstititor)

	ss.CreateVariable("foo", "123")
	ss.CreateVariable("food", "456")

	testExpansion(t, ss, " foo 456", " 123 456")
	testExpansion(t, ss, " foo food", " 123 456")
	testExpansion(t, ss, "", "")
	testExpansion(t, ss, " ", " ")
	testExpansion(t, ss, "1", "1")
}

func testExpansion(t *testing.T, ss StrSubstititor, in string, expected string) {
	out := ss.Expand(in)
	if out != expected {
		t.Errorf("Expanding \"%s\":\nReceived: \"%s\"\nExpected: \"%s\"", in, out, expected)
	}
}
