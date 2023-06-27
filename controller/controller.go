package controller

import (
	"fmt"
	"sync"
	d "ticket-pimp/domain"
	"ticket-pimp/ext"
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
}

func (wc *WorkflowController) Workflow(name string) (string, error) {
	yt := wc.iYouTrack

	projects, err := yt.GetProjects()

	if err != nil {
		return "", err
	}

	issue, err := yt.CreateIssue(projects[1].ID, name, "")

	if err != nil {
		return "", err
	}

	if issue != nil {
		var (
			git, gitBuild *d.Git
			cloud         *d.Folder
		)

		var wg sync.WaitGroup
		wg.Add(3)

		go func(ref **d.Git) {
			defer wg.Done()
			*ref, _ = wc.CreateRepo(issue.Key)
		}(&git)

		go func(ref **d.Git) {
			defer wg.Done()
			*ref, _ = wc.CreateRepo(issue.Key + "-build")
		}(&gitBuild)

		go func(ref **d.Folder) {
			defer wg.Done()
			*ref, _ = wc.CreateFolder(issue.Key + " - " + issue.Summary)
		}(&cloud)

		wg.Wait()

		yt.UpdateIssue(
			issue,
			cloud.PrivateURL,
			git.HtmlUrl,
			fmt.Sprintf("ssh://%s/%s.git", gitBuild.SshUrl, gitBuild.FullName))
	}
	return issue.Key, nil
}
