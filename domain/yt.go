package domain

import (
	"log"
	"time"

	"github.com/imroc/req/v3"
)

type youtrack struct {
	client *req.Client
}

func NewYT(base, token string) *youtrack {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	client := req.C().
		SetTimeout(15 * time.Second).
		SetCommonHeaders(headers).
		SetBaseURL(base).
		SetCommonBearerAuthToken(token)

	return &youtrack{
		client: client,
	}
}

type Project struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}

// GetProjects
// provides an array of existing projects;
func (yt *youtrack) GetProjects() []Project {

	var projects []Project

	resp, err := yt.client.R().
		EnableDump().
		SetQueryParam("fields", "id,name,shortName").
		SetSuccessResult(&projects).
		Get("/admin/projects")

	if !resp.IsSuccessState() || err != nil {
		log.Print("bad status:", resp.Status)
		log.Print(resp.Dump())
	}
	return projects
}

type ProjectID struct {
	ID string `json:"id"`
}

type IssueCreateRequest struct {
	ProjectID   ProjectID `json:"project"`
	Key         string    `json:"idReadable"`
	ID          string    `json:"id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
}

// CreateIssue
// example: newIssue := yt.CreateIssue("0-2", "Summary", "Description");
func (yt *youtrack) CreateIssue(projectID, name string) *IssueCreateRequest {

	// Create an issue with the provided:, Project ID, Name, Description;
	issue := IssueCreateRequest{
		ProjectID: ProjectID{
			ID: projectID, //"id":"0-2"
		},
		Summary: name,
		//Description: description,
	}

	// Push issue to the YT;
	resp, err := yt.client.R().
		SetQueryParam("fields", "idReadable,id").
		SetBody(&issue).
		SetSuccessResult(&issue).
		Post("/issues")

	// Check if request failed or response status is not Ok;
	if !resp.IsSuccessState() || err != nil {
		log.Print("bad status:", resp.Status)
		log.Print(resp.Dump())
	}

	return &issue
}

type IssueUpdateRequest struct {
	IssueCreateRequest
	CustomFields []CustomField `json:"customFields"`
}

type CustomFields struct {
	List []CustomField `json:"customFields"`
}

type CustomField struct {
	Name  string `json:"name"`
	Type  string `json:"$type"`
	Value string `json:"value"`
}

func (yt *youtrack) UpdateIssue(issue *IssueCreateRequest, folder, git, gitBuild string) *IssueUpdateRequest {
	// Set Folder, Git, GitBuild to the Issue:
	update := IssueUpdateRequest{
		IssueCreateRequest: *issue,
		CustomFields: []CustomField{
			{
				Name:  "Директория графики",
				Type:  "SimpleIssueCustomField",
				Value: folder,
			},
			{
				Name:  "Репо проекта",
				Type:  "SimpleIssueCustomField",
				Value: git,
			},
			{
				Name:  "Репо iOS сборки",
				Type:  "SimpleIssueCustomField",
				Value: gitBuild,
			},
		},
	}

	// Push issue update to  YT
	resp, err := yt.client.R().
		SetBody(&update).
		SetSuccessResult(&issue).
		Post("/issues/" + issue.Key)

	if !resp.IsSuccessState() || err != nil {
		log.Print("bad status:", resp.Status)
		log.Print(resp.Dump())
	}

	return &update
}
