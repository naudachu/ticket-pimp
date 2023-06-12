package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"ticket-creator/controller"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

func errorAnswer(errorMsg string) string {
	return tg.HTML.Text(
		tg.HTML.Line(
			tg.HTML.Italic(errorMsg),
		),
	)
}

func NewTicketHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	str := strings.Replace(mu.Text, "/new", "", 1)

	if str == "" {
		return errors.New("empty command provided")
	}

	issueKeyStr, err := controller.Workflow(str)

	if err != nil {
		return mu.Answer(errorAnswer(err.Error())).ParseMode(tg.HTML).DoVoid(ctx)
	}

	return mu.Answer(newTicketAnswer(issueKeyStr)).ParseMode(tg.HTML).DoVoid(ctx)
}

func newTicketAnswer(name string) string {
	return tg.HTML.Text(
		tg.HTML.Line(
			"🤘 Ticket ",
			tg.HTML.Link(name, fmt.Sprintf("https://marlerino.youtrack.cloud/issue/%s", name)),
			"has been created!",
		),
	)
}

func NewRepoHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	str := strings.Replace(mu.Text, "/repo", "", 1)

	if str == "" {
		return errors.New("empty command provided")
	}

	repoStr, err := controller.CreateRepo(str, 0)

	if err != nil {
		return mu.Answer(errorAnswer(err.Error())).ParseMode(tg.HTML).DoVoid(ctx)
	}

	return mu.Answer(newRepoAnswer(repoStr)).ParseMode(tg.HTML).DoVoid(ctx)
}

func newRepoAnswer(name string) string {
	return tg.HTML.Text(
		tg.HTML.Line(
			"Repo ",
			name,
			"has been created!",
		),
	)
}

func PingHandler(ctx context.Context, mu *tgb.MessageUpdate) error {
	return mu.Answer("pong").DoVoid(ctx)
}
