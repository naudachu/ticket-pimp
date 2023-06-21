package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"ticket-pimp/controller"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

type Handler struct {
	workflow controller.IWorkflowController
}

func NewHandler(gitBaseURL, gitToken, cloudBaseURL, cloudAuthUser, cloudAuthPass, ytBaseURL, ytToken string) *Handler {
	return &Handler{
		workflow: controller.NewWorkflowController(
			gitBaseURL,
			gitToken,
			cloudBaseURL,
			cloudAuthUser,
			cloudAuthPass,
			ytBaseURL,
			ytToken),
	}
}

func (h *Handler) PingHandler(ctx context.Context, mu *tgb.MessageUpdate) error {
	return mu.Answer("pong").DoVoid(ctx)
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

func (h *Handler) NewRepoHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	str := strings.Replace(mu.Text, "/repo", "", 1)

	if str == "" {
		return errors.New("empty command provided")
	}

	repoStr, err := h.workflow.CreateRepo(str, 0)

	if err != nil {
		return mu.Answer(errorAnswer(err.Error())).ParseMode(tg.HTML).DoVoid(ctx)
	}

	return mu.Answer(newRepoAnswer(repoStr)).ParseMode(tg.HTML).DoVoid(ctx)
}

func (h *Handler) NewFolderHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	str := strings.Replace(mu.Text, "/folder", "", 1)

	if str == "" {
		return errors.New("empty command provided")
	}

	cloud, err := h.workflow.CreateFolder(str)

	if err != nil {
		return mu.Answer(errorAnswer(err.Error())).ParseMode(tg.HTML).DoVoid(ctx)
	}

	answer := tg.HTML.Text(
		tg.HTML.Line(
			"✨ Shiny folder",
			tg.HTML.Link(cloud.Title, cloud.PrivateURL),
			"has been created!",
		),
	)

	return mu.Answer(answer).
		ParseMode(tg.HTML).
		DoVoid(ctx)
}

func errorAnswer(errorMsg string) string {
	return tg.HTML.Text(
		tg.HTML.Line(
			tg.HTML.Italic(errorMsg),
		),
	)
}

func (h *Handler) NewTicketHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	str := strings.Replace(mu.Text, "/new", "", 1)

	if str == "" {
		return errors.New("empty command provided")
	}

	issueKeyStr, err := h.workflow.Workflow(str)

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
