package workflow

import (
	"os"
	"ticket-creator/domain"
)

func ProduceTicket(name string) (string, error) {
	yt := domain.NewYT(os.Getenv("YT_URL"), os.Getenv("YT_TOKEN"))

	projects, err := yt.GetProjects()
	if err != nil {
		return "", err
	}

	issue, err := yt.CreateIssue(projects[1].ID, name)
	if err != nil {
		return "", err
	}

}
