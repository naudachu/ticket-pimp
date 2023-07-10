package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"ticket-pimp/internal/controllers"
	"ticket-pimp/internal/controllers/controller"
	"ticket-pimp/internal/storage"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

type Handler struct {
	workflow controllers.IWorkflowController
	cloud    controller.CloudCreator
	git      controller.RepoCreator
}

func NewHandler(git controller.RepoCreator, cloud controller.CloudCreator, workflow controllers.IWorkflowController, r storage.Storage) *Handler {
	return &Handler{
		workflow: workflow,
		cloud:    cloud,
		git:      git,
	}
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

func (h *Handler) NewTaskHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	taskText := strings.TrimSpace(strings.Replace(mu.Text, "/task", "", 1))
	words := strings.Split(taskText, " ")

	var summaryTail string
	if len(words) > 3 {
		summaryTail = strings.Join(words[0:3], " ")
	} else {
		summaryTail = strings.Join(words, " ")
	}

	task := h.workflow.NewTask(
		summaryTail,
		taskText,
		mu.From.Username.PeerID(),
		mu.From.Username.Link(),
	)

	createdTicket, err := h.workflow.CreateTask(task)
	if err != nil {
		return mu.Answer(errorAnswer(err.Error())).ParseMode(tg.HTML).DoVoid(ctx)
	}

	return mu.Answer(tg.HTML.Text(
		tg.HTML.Line(
			"🤘 Задача",
			tg.HTML.Link(createdTicket.Key, createdTicket.URL),
			"была создана!",
		),
	)).ParseMode(tg.HTML).DoVoid(ctx)
}

func (h *Handler) NewConversion(ctx context.Context, mu *tgb.MessageUpdate) error {
	msg := strings.TrimSpace(strings.Replace(mu.Caption, "/conversion", "", 1))

	appID, token := normalizeToken(msg)

	fid := mu.Update.Message.Document.FileID

	client := mu.Client

	file, err := client.GetFile(fid).Do(ctx)
	if err != nil {
		return err
	}

	f, err := client.Download(ctx, file.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	l := h.workflow.ThrowConversions(f, appID, token)

	if len(l.Advertiser) != 0 {
		return mu.Answer(tg.HTML.Text(
			"Неуспешные запросы:",
			tg.HTML.Code(strings.Join(l.Advertiser, ", ")),
		)).ParseMode(tg.HTML).DoVoid(ctx)
	}

	return mu.Answer(tg.HTML.Text(
		"Конверсии отправлены",
	)).ParseMode(tg.HTML).DoVoid(ctx)
}

func normalizeToken(msg string) (string, string) {
	msg = strings.TrimSpace(msg)

	args := strings.Split(msg, "|")

	if len(args) != 2 {
		log.Print(len(args))
		return "", ""
	}

	return args[0], args[0] + "|" + args[1]
}
