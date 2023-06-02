package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ticket-creator/domain"

	"github.com/joho/godotenv"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

func main() {
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

func run(ctx context.Context) error {
	client := tg.New(os.Getenv("TG_API"))

	router := tgb.NewRouter().
		Message(func(ctx context.Context, mu *tgb.MessageUpdate) error {

			str := strings.Replace(mu.Text, "/new", "", 1)
			if str == "" {
				return errors.New("empty command provided")
			}
			issueKeyStr := workflow(str)

			return mu.Answer(tg.HTML.Text(
				tg.HTML.Line(
					"🤘 Ticket ",
					tg.HTML.Link(issueKeyStr, fmt.Sprintf("https://marlerino.youtrack.cloud/issue/%s", issueKeyStr)),
					"has been created!",
				),
			)).ParseMode(tg.HTML).DoVoid(ctx)
		}, tgb.TextHasPrefix("/new"))

	return tgb.NewPoller(
		router,
		client,
	).Run(ctx)
}

func workflow(name string) string {
	yt := domain.NewYT(os.Getenv("YT_URL"), os.Getenv("YT_TOKEN"))
	projects := yt.GetProjects()
	issue := yt.CreateIssue(projects[0].ID, name)
	if issue != nil {
		git := createRepo(issue.Key, 0)
		gitBuild := createRepo(issue.Key+"-build", 1)
		folder := createFolder(issue.Key + " - " + issue.Summary)
		yt.UpdateIssue(issue, folder, git, gitBuild)
	}
	return issue.Key
}

func createRepo(name string, param uint) string {
	gb := domain.NewGitBucket(os.Getenv("GIT_BASE_URL"), os.Getenv("GIT_TOKEN"))
	repo, _ := gb.NewRepo(name)
	if repo != nil {
		switch param {
		case 0:
			return repo.HtmlUrl
		case 1:
			return fmt.Sprintf("ssh://%s/%s.git", repo.SshUrl, repo.FullName)
		default:
			return repo.CloneUrl
		}
	}
	return "no-repo"
}

func createFolder(name string) string {
	oc := domain.NewCloud(os.Getenv("CLOUD_BASE_URL"), os.Getenv("CLOUD_USER"), os.Getenv("CLOUD_PASS"))
	cloud, _ := oc.CreateFolder(name)
	if cloud != nil {
		return cloud.FolderPath
	}
	return "no-folder"
}
