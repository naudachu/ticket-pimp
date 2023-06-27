package controller

import "fmt"

type Task struct {
	Summary     string
	Description string
	Creator     string
	CreatorLink string

	Key string
	URL string
}

func (wc *WorkflowController) NewTask(summ, desc, c, cLink string) *Task {
	return &Task{
		Summary:     summ,
		Description: desc,
		Creator:     c,
		CreatorLink: cLink,
	}
}

func (wc *WorkflowController) CreateTask(t *Task) (*Task, error) {

	yt := wc.iYouTrack

	projects, err := yt.GetProjects()
	if err != nil {
		return nil, err
	}

	t.Description += fmt.Sprintf("\n\n Created by: [%s](%s)", t.Creator, t.CreatorLink)

	issue, err := yt.CreateIssue(projects[1].ID, t.Creator+" | "+t.Summary, t.Description)
	if err != nil {
		return nil, err
	}

	t.Key = issue.Key
	t.URL = fmt.Sprintf("https://marlerino.youtrack.cloud/issue/%s", issue.Key)

	return t, nil //[ ] normal return;
}
