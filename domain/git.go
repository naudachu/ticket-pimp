package domain

import (
	"log"
	"time"

	"github.com/imroc/req/v3"
)

type gitbucket struct {
	client *req.Client
}

func NewGitBucket(base, token string) *gitbucket {
	headers := map[string]string{
		"Accept":               "application/vnd.github+json",
		"Authorization":        "Token " + token,
		"X-GitHub-Api-Version": "2022-11-28",
		"Content-Type":         "application/json",
	}

	client := req.C().
		SetTimeout(5 * time.Second).
		SetCommonHeaders(headers).
		SetBaseURL(base)
	return &gitbucket{
		client: client,
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

	resp, err := gb.client.R().
		SetBody(&payload).
		SetSuccessResult(&git).
		Post("/user/repos")

	// Check if request failed or response status is not Ok;
	if !resp.IsSuccessState() || err != nil {
		log.Print("bad status:", resp.Status)
		log.Print(resp.Dump())
	}

	return &git, err
}
