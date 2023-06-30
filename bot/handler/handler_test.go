package handler

import (
	"testing"
	"ticket-pimp/bot/domain"
)

type test struct {
	arg      domain.Git
	expected string
}

var tests = []test{
	{domain.Git{
		Name:     "text",
		FullName: "",
		Private:  false,
		Url:      "",
		CloneUrl: "",
		HtmlUrl:  "https://reddit.com/",
		SshUrl:   "",
	}, "Repo  <a href=\"https://reddit.com/\">text</a> has been created!"},
}

func TestPrepareAnswer(t *testing.T) {

	for _, test := range tests {
		g := newGit(&test.arg)

		if output := g.PrepareAnswer(); output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}
