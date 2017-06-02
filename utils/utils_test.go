package utils

import (
	"testing"
)

func TestStrSubstitutor(t *testing.T) {
	ss := make(StrSubstititor)

	ss.CreateVariable("foo", "123")
	ss.CreateVariable("food", "456")

	testExpansion(t, ss, " $foo 456", " 123 456")
	testExpansion(t, ss, " $foo $food", " 123 456")
	testExpansion(t, ss, " foo 456", " foo 456")
	testExpansionError(t, ss, "$")
	testExpansionError(t, ss, "$bar")
	testExpansionError(t, ss, "foo$foo")
	testExpansionError(t, ss, "$fool")
}

func testExpansion(t *testing.T, ss StrSubstititor, in string, expected string) {
	out, err := ss.Expand(in)
	if err != nil {
		t.Errorf("Expanding \"%s\":\nError: %s", in, err)
		return
	}
	if out != expected {
		t.Errorf("Expanding \"%s\":\nReceived: \"%s\"\nExpected: \"%s\"", in, out, expected)
	}
}

func testExpansionError(t *testing.T, ss StrSubstititor, in string) {
	_, err := ss.Expand(in)
	if err == nil {
		t.Errorf("Expanding \"%s\":\nExpected to receive an error", in)
	}
}
