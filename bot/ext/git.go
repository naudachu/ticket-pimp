package ext

import (
	"log"
	"os"
	"ticket-pimp/bot/domain"
	"ticket-pimp/bot/helpers"
	"time"
)

type Git struct {
	*Client
}

type IGit interface {
	NewRepo(string) (*domain.Git, error)
	AppsAsCollaboratorTo(*domain.Git) (*domain.Git, error)
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
	git.Private = true

	resp, _ := gb.R().
		SetBody(&payload).
		SetSuccessResult(&git).
		Post("/user/repos")

	if resp.Err != nil {
		log.Print(resp.Err)
		return nil, resp.Err
	}

	return &git, nil
}

func (gb *Git) AppsAsCollaboratorTo(git *domain.Git) (*domain.Git, error) {

	payload := permissionRequest{
		Perm: "admin",
	}

	respURL := "/repos/" + os.Getenv("GIT_USER") + "/" + git.Name + "/collaborators/apps"

	resp, _ := gb.R().
		SetBody(&payload).
		Put(respURL)

	if resp.Err != nil {
		log.Print(resp.Err)
		return nil, resp.Err
	}

	return git, nil
}
