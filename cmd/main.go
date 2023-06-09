package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"ticket-creator/domain"

	"github.com/joho/godotenv"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

func main() {
	log.Print("started")
	env(".env")
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Println(err)
		defer os.Exit(1)
	}
}

func env(envFilePath string) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func errorAnswer(errorMsg string) string {
	return tg.HTML.Text(
		tg.HTML.Line(
			tg.HTML.Italic(errorMsg),
		),
	)
}

func run(ctx context.Context) error {
	client := tg.New(os.Getenv("TG_API"))

	router := tgb.NewRouter().
		Message(newTicketHandler, tgb.TextHasPrefix("/new")).
		Message(pingHandler, tgb.Command("ping")).
		Message(newRepoHandler, tgb.TextHasPrefix("/repo"))

	return tgb.NewPoller(
		router,
		client,
	).Run(ctx)
}

func newTicketHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	str := strings.Replace(mu.Text, "/new", "", 1)
	if str == "" {
		return errors.New("empty command provided")
	}
	issueKeyStr, err := workflow(str)
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

func newRepoHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	str := strings.Replace(mu.Text, "/newrepo", "", 1)
	if str == "" {
		return errors.New("empty command provided")
	}
	repoStr, err := createRepo(str, 0)
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

func pingHandler(ctx context.Context, mu *tgb.MessageUpdate) error {
	return mu.Answer("pong").DoVoid(ctx)
}

func workflow(name string) (string, error) {
	yt := domain.NewYT(os.Getenv("YT_URL"), os.Getenv("YT_TOKEN"))

	projects, err := yt.GetProjects()
	if err != nil {
		return "", err
	}

	issue, err := yt.CreateIssue(projects[1].ID, name)
	if err != nil {
		return "", err
	}
	if issue != nil {
		var (
			wg                    sync.WaitGroup
			git, gitBuild, folder string
		)

		wg.Add(3)

		go func() {
			defer wg.Done()
			git, _ = createRepo(issue.Key, 0)
		}()

		go func() {
			defer wg.Done()
			gitBuild, _ = createRepo(issue.Key+"-build", 1)
		}()
		go func() {
			defer wg.Done()
			folder = createFolder(issue.Key + " - " + issue.Summary)
		}()

		wg.Wait()

		yt.UpdateIssue(issue, folder, git, gitBuild)
	}
	return issue.Key, nil
}

func createRepo(name string, param uint) (string, error) {
	gb := domain.NewGitBucket(os.Getenv("GIT_BASE_URL"), os.Getenv("GIT_TOKEN"))
	repo, err := gb.NewRepo(name)
	if err != nil {
		return "", err
	}

	// Result string formatting:
	if repo != nil {
		switch param {
		case 0:
			return repo.HtmlUrl, nil
		case 1:
			return fmt.Sprintf("ssh://%s/%s.git", repo.SshUrl, repo.FullName), nil
		default:
			return repo.CloneUrl, nil
		}
	}

	return "", err
}

func createFolder(name string) string {
	oc := domain.NewCloud(os.Getenv("CLOUD_BASE_URL"), os.Getenv("CLOUD_USER"), os.Getenv("CLOUD_PASS"))
	cloud, _ := oc.CreateFolder(name)
	if cloud != nil {
		return cloud.FolderPath
	}
	return "no-folder"
}
