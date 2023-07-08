package domain

import "fmt"

type Project struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}

type ProjectID struct {
	ID string `json:"id"`
}

// Find needed project.ID in the project's list
func (plist *ProjectsList) FindProjectByName(searchName string) (string, error) {

	projectID := ""

	for _, elem := range plist.Projects {
		if elem.ShortName == searchName {
			projectID = elem.ID
		}
	}

	if projectID == "" {
		return "", fmt.Errorf("project %s doesn't exist", searchName)
	}
	return projectID, nil
}

type IssueCreateRequest struct {
	ProjectID   ProjectID `json:"project"`
	Key         string    `json:"idReadable"`
	ID          string    `json:"id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
}

// [ ] try `,omitempty` to remove extra struct;

type IssueUpdateRequest struct {
	IssueCreateRequest
	CustomFields []CustomField `json:"customFields"`
}

type CustomField struct {
	Name  string `json:"name"`
	Type  string `json:"$type"`
	Value string `json:"value"`
}

type ProjectsList struct {
	Projects []Project
}
