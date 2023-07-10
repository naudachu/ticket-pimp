package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"ticket-pimp/internal/domain"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

type git struct {
	name string
	url  string

	git string
	ssh string
}

func newGit(d *domain.Git) *git {
	return &git{
		name: d.Name,
		url:  d.HtmlUrl,
		git:  d.CloneUrl,
		ssh:  fmt.Sprintf("ssh://%s/%s.git", d.SshUrl, d.FullName),
	}
}

// FYI: Telegram doesn't renders this hyperlink, if the url is localhost 🤷‍♂️
func (g *git) PrepareAnswer() string {
	return tg.HTML.Text(
		tg.HTML.Line(
			"Repo ",
			tg.HTML.Link(g.name, g.url),
			"has been created!",
		),
	)
}

func (h *Handler) NewRepoHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	str := strings.Replace(mu.Text, "/repo", "", 1)

	if str == "" {
		return errors.New("empty command provided")
	}

	var g *domain.Git
	g, err := h.git.CreateRepo(str)

	if err != nil {
		return mu.Answer(errorAnswer(err.Error())).ParseMode(tg.HTML).DoVoid(ctx)
	}

	resp := newGit(g).PrepareAnswer()

	return mu.Answer(resp).ParseMode(tg.HTML).DoVoid(ctx)
}
