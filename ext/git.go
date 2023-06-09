package ext

import (
	"log"
	"regexp"
	"strings"
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

func gitHubLikeNaming(input string) string {
	// Remove leading and trailing whitespace
	input = strings.TrimSpace(input)

	// Replace non-Latin letters with spaces
	reg := regexp.MustCompile("[^a-zA-Z]+")
	input = reg.ReplaceAllString(input, " ")

	// Split into words and capitalize first letter of each
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	// Join words and return
	return strings.Join(words, "-")
}

func (gb *gitbucket) NewRepo(name string) (*Repo, error) {
	name = gitHubLikeNaming(name)

	type request struct {
		Name    string `json:"name"`
		Private bool   `json:"private"`
	}

	payload := request{
		Name:    name,
		Private: false,
	}

	var git Repo

	resp, err := gb.R().
		SetBody(&payload).
		SetSuccessResult(&git).
		Post("/user/repos")

	if err != nil {
		log.Print(resp)
		return nil, err
	}

	type permissionRequest struct {
		Perm string `json:"permission"`
	}

	payloadPermission := permissionRequest{
		Perm: "admin",
	}

	resp, err = gb.R().
		SetBody(&payloadPermission).
		Put("/repos/naudachu/" + name + "/collaborators/apps")

	if err != nil {
		log.Print(resp)
		return nil, err
	}

	return &git, err
}
