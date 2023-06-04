package domain

import (
	"time"

	"github.com/imroc/req/v3"
)

type gitbucket struct {
	*req.Client
}

func NewGitBucket(base, token string) *gitbucket {
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
	return &gitbucket{
		client,
	}
}

type Repo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Url      string `json:"url"`
	CloneUrl string `json:"clone_url"`
	HtmlUrl  string `json:"Html_url"`
	SshUrl   string `json:"ssh_url"`
}

func (gb *gitbucket) NewRepo(name string) (*Repo, error) {

	type request struct {
		Name    string `json:"name"`
		Private bool   `json:"private"`
	}

	payload := request{
		Name:    name,
		Private: false,
	}

	var git Repo

	_, err := gb.R().
		SetBody(&payload).
		SetSuccessResult(&git).
		Post("/user/repos")

	if err != nil {
		return nil, err
	}

	type permissionRequest struct {
		Perm string `json:"permission"`
	}

	payloadPermission := permissionRequest{
		Perm: "admin",
	}

	_, err = gb.R().
		SetBody(&payloadPermission).
		Put("/repos/naudachu/" + name + "/collaborators/apps")

	if err != nil {
		return nil, err
	}

	return &git, err
}
