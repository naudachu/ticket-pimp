package helpers

import (
	"testing"
)

type test struct {
	arg, expected string
}

var tests = []test{
	{" App-21", "app-21"},
	{"App-21-build", "app-21-build"},
	{" hello - biatch", "hello-biatch"},
	{"  `~!@#$%^&*()=+ abc `~!@#$%^&*()=+ 2-22 ", "abc-2-22"},
}

func TestGitNaming(t *testing.T) {

	for _, test := range tests {
		if output := GitNaming(test.arg); output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}
