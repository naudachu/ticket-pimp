package ext

import (
	"log"
	"os"
	"ticket-creator/domain"
	"ticket-creator/helpers"
	"time"
)

type Git struct {
	*Client
	*domain.Git
}

func NewGit(base, token string) *Git {
	headers := map[string]string{
		"Accept":               "application/vnd.github+json",
		"Authorization":        "Token " + token,
		"X-GitHub-Api-Version": "2022-11-28",
		"Content-Type":         "application/json",
	}

	client := NewClient().
		SetTimeout(5 * time.Second).
		SetCommonHeaders(headers).
		SetBaseURL(base)

	return &Git{
		Client: &Client{client},
		Git: &domain.Git{
			Name:     "",
			FullName: "",
			Private:  true,
			Url:      "",
			CloneUrl: "",
			HtmlUrl:  "",
			SshUrl:   "",
		},
	}
}

type request struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

type permissionRequest struct {
	Perm string `json:"permission"`
}

func (gb *Git) NewRepo(name string) (*domain.Git, error) {
	name = helpers.GitNaming(name)

	payload := request{
		Name:    name,
		Private: true,
	}

	var git domain.Git

	resp, err := gb.R().
		SetBody(&payload).
		SetSuccessResult(&git).
		Post("/user/repos")
		//Post("/orgs/apps/repos")

	if err != nil {
		log.Print(resp)
	}

	return &git, err
}

func (gb *Client) AppsAsCollaboratorTo(git *domain.Git) (*domain.Git, error) {
	payloadPermission := permissionRequest{
		Perm: "admin",
	}

	resp, err := gb.R().
		SetBody(&payloadPermission).
		Put("/repos/" + os.Getenv("GIT_USER") + "/" + git.Name + "/collaborators/apps")

	if err != nil {
		log.Print(resp)
	}

	return git, err
}
