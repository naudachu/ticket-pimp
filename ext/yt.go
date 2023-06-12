package ext

import (
	"fmt"
	"log"
	"time"

	"github.com/imroc/req/v3"
)

type youtrack struct {
	*req.Client
}

func NewYT(base, token string) *youtrack {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	client := NewClient().
		SetTimeout(15 * time.Second).
		SetCommonHeaders(headers).
		SetBaseURL(base).
		SetCommonBearerAuthToken(token)

	return &youtrack{
		client,
	}
}

type Project struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}

// GetProjects
// provides an array of existing projects;
func (yt *youtrack) GetProjects() ([]Project, error) {

	var projects []Project

	_, err := yt.R().
		EnableDump().
		SetQueryParam("fields", "id,name,shortName").
		SetSuccessResult(&projects).
		Get("/admin/projects")

	// Check if the request failed;
	if err != nil {
		return nil, fmt.Errorf("some problem with YT request. error message: %v", err)
	}

	return projects, nil
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
func (yt *youtrack) CreateIssue(projectID, name string) (*IssueCreateRequest, error) {

	// Create an issue with the provided:, Project ID, Name, Description;
	issue := IssueCreateRequest{
		ProjectID: ProjectID{
			ID: projectID, //"id":"0-2"
		},
		Summary: name,
		//Description: description,
	}

	// Push issue to the YT;
	_, err := yt.R().
		SetQueryParam("fields", "idReadable,id").
		SetBody(&issue).
		SetSuccessResult(&issue).
		Post("/issues")

	// Check if the request failed;
	if err != nil {
		return nil, fmt.Errorf("some problem with YT request. error message: %v", err)
	}

	return &issue, nil
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

func (yt *youtrack) UpdateIssue(issue *IssueCreateRequest, folder, git, gitBuild string) (*IssueUpdateRequest, error) {
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
	resp, err := yt.R().
		SetBody(&update).
		SetSuccessResult(&issue).
		Post("/issues/" + issue.Key)

		// Check if the request failed;
	if err != nil {
		return nil, fmt.Errorf("some problem with YT request. error message: %v", err)
	}

	if !resp.IsSuccessState() {
		log.Print("bad status:", resp.Status)
		log.Print(resp.Dump())
		return nil, fmt.Errorf("YouTrack responded with %d", resp.StatusCode)
	}

	return &update, nil
}
