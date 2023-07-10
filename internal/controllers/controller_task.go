package controllers

import (
	"fmt"
	"ticket-pimp/internal/domain"

	"gorm.io/gorm"
)

type YTTask struct {
	Summary     string
	Description string
	Creator     string
	CreatorLink string

	Key string
	URL string
}

func (wc *WorkflowController) NewTask(summ, desc, c, cLink string) *YTTask {
	return &YTTask{
		Summary:     summ,
		Description: desc,
		Creator:     c,
		CreatorLink: cLink,
	}
}

func (wc *WorkflowController) CreateTask(t *YTTask) (*YTTask, error) {

	yt := wc.iYouTrack

	projectID, err := yt.GetProjectIDByName("tst")
	if err != nil {
		return nil, err
	}

	t.Description += fmt.Sprintf("\n\n Created by: [%s](%s)", t.Creator, t.CreatorLink)

	task := domain.TaskEntity{
		Title:       t.Summary,
		Description: t.Description,
		Creator:     t.Creator,
		Responsible: "",
		Status:      0,
		Model:       gorm.Model{},
	}

	taskFromDB, err := wc.taskRepo.SaveTask(&task)
	if err != nil {
		return nil, err //[ ] переделать
	}

	fmt.Print(taskFromDB)

	issue, err := yt.CreateIssue(projectID, t.Creator+" | "+t.Summary, t.Description)
	if err != nil {
		return nil, err
	}

	t.Key = issue.Key
	t.URL = fmt.Sprintf("https://marlerino.youtrack.cloud/issue/%s", issue.Key)

	return t, nil
}
