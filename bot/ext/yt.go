package ext

import (
	"fmt"
	"log"
	"time"

	d "ticket-pimp/bot/domain"

	"github.com/imroc/req/v3"
)

type youtrack struct {
	*req.Client
}

type IYouTrack interface {
	GetProjectIDByName(searchName string) (string, error)
	CreateIssue(projectID, name, description string) (*d.IssueCreateRequest, error)
	UpdateIssue(issue *d.IssueCreateRequest, folder, git, gitBuild string) (*d.IssueUpdateRequest, error)
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

// GetProjects
// provides an array of existing projects;
func (yt *youtrack) getProjects() (*d.ProjectsList, error) {

	var projects d.ProjectsList

	resp, _ := yt.R().
		EnableDump().
		SetQueryParam("fields", "id,name,shortName").
		SetSuccessResult(&projects.Projects).
		Get("/admin/projects")

	// Check if the request failed;
	if resp.Err != nil {
		return nil, fmt.Errorf("some problem with YT request. error message: %v", resp.Err)
	}

	return &projects, nil
}

// GetProjects
// provides an array of existing projects;
func (yt *youtrack) GetProjectIDByName(searchName string) (string, error) {

	projects, err := yt.getProjects()
	if err != nil {
		return "", err
	}

	projectID, err := projects.FindProjectByName(searchName)
	if err != nil {
		return "", err
	}

	return projectID, nil
}

// CreateIssue
// example: newIssue := yt.CreateIssue("0-2", "Summary", "Description");
func (yt *youtrack) CreateIssue(projectID, name string, description string) (*d.IssueCreateRequest, error) {

	// Create an issue with the provided:, Project ID, Name, Description;
	issue := d.IssueCreateRequest{
		ProjectID: d.ProjectID{
			ID: projectID, //"id":"0-2"
		},
		Summary:     name,
		Description: description,
	}

	// Push issue to the YT;
	resp, _ := yt.R().
		SetQueryParam("fields", "idReadable,id").
		SetBody(&issue).
		SetSuccessResult(&issue).
		Post("/issues")

	// Check if the request failed;
	if resp.Err != nil {
		return nil, fmt.Errorf("some problem with YT request. error message: %v", resp.Err)
	}

	return &issue, nil
}

func (yt *youtrack) UpdateIssue(issue *d.IssueCreateRequest, folder, git, gitBuild string) (*d.IssueUpdateRequest, error) {
	// Set Folder, Git, GitBuild to the Issue:
	update := d.IssueUpdateRequest{
		IssueCreateRequest: *issue,
		CustomFields: []d.CustomField{
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
	resp, _ := yt.R().
		SetBody(&update).
		SetSuccessResult(&issue).
		Post("/issues/" + issue.Key)

		// Check if the request failed;
	if resp.Err != nil {
		return nil, fmt.Errorf("some problem with YT request. error message: %v", resp.Err)
	}

	if !resp.IsSuccessState() {
		log.Print("bad status:", resp.Status)
		log.Print(resp.Dump())
		return nil, fmt.Errorf("YouTrack responded with %d", resp.StatusCode)
	}

	return &update, nil
}
