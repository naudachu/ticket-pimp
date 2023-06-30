package controller

import (
	"io"
	d "ticket-pimp/bot/domain"
	"ticket-pimp/bot/ext"
)

type WorkflowController struct {
	iGit      ext.IGit
	iCloud    ext.ICloud
	iYouTrack ext.IYouTrack
	iCoda     ext.ICoda
}

func NewWorkflowController(
	gitBaseURL,
	gitToken,
	cloudBaseURL,
	cloudAuthUser,
	cloudAuthPass,
	ytBaseURL,
	ytToken string,
) *WorkflowController {
	return &WorkflowController{
		iGit:      ext.NewGit(gitBaseURL, gitToken),
		iCloud:    ext.NewCloud(cloudBaseURL, cloudAuthUser, cloudAuthPass),
		iYouTrack: ext.NewYT(ytBaseURL, ytToken),
		iCoda:     ext.NewCodaClient(),
	}
}

type IWorkflowController interface {
	Workflow(name string) (string, error)
	CreateRepo(name string) (*d.Git, error)
	CreateFolder(name string) (*d.Folder, error)

	NewTask(summ, desc, c, cLink string) *Task
	CreateTask(t *Task) (*Task, error)

	ThrowConversions(f io.ReadCloser, appID string, token string) *d.ConversionLog
}
