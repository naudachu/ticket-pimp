package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"ticket-pimp/bot/controller"
	d "ticket-pimp/bot/domain"

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

type git struct {
	name string
	url  string

	git string
	ssh string
}

func newGit(d *d.Git) *git {
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

	var g *d.Git
	g, err := h.workflow.CreateRepo(str)

	if err != nil {
		return mu.Answer(errorAnswer(err.Error())).ParseMode(tg.HTML).DoVoid(ctx)
	}

	resp := newGit(g).PrepareAnswer()

	return mu.Answer(resp).ParseMode(tg.HTML).DoVoid(ctx)
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
