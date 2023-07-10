package controllers

import (
	"io"
	"ticket-pimp/internal/controllers/controller"
	d "ticket-pimp/internal/domain"
	"ticket-pimp/internal/extapi"
	"ticket-pimp/internal/storage"
)

type WorkflowController struct {
	git   controller.RepoCreator
	cloud controller.CloudCreator

	iYouTrack extapi.IYouTrack
	iCoda     extapi.ICoda
	taskRepo  storage.Storage
}

func NewWorkflowController(
	ytBaseURL,
	ytToken string,

	r storage.Storage,
	git controller.RepoCreator,
	cloud controller.CloudCreator,

) *WorkflowController {
	return &WorkflowController{
		git:       git,
		cloud:     cloud,
		iYouTrack: extapi.NewYT(ytBaseURL, ytToken),
		iCoda:     extapi.NewCodaClient(),
		taskRepo:  r,
	}
}

type IWorkflowController interface {
	Workflow(name string) (string, error)
	NewTask(summ, desc, c, cLink string) *YTTask
	CreateTask(t *YTTask) (*YTTask, error)
	ThrowConversions(f io.ReadCloser, appID string, token string) *d.ConversionLog
}
