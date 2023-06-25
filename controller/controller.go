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
	}
}

type IWorkflowController interface {
	Workflow(name string) (string, error)
	CreateRepo(name string) (*d.Git, error)
	CreateFolder(name string) (*d.Folder, error)
}

func (wc *WorkflowController) Workflow(name string) (string, error) {
	yt := wc.iYouTrack

	projects, err := yt.GetProjects()

	if err != nil {
		return "", err
	}

	issue, err := yt.CreateIssue(projects[1].ID, name)

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

func (wc *WorkflowController) CreateRepo(name string) (*d.Git, error) {
	//Create git repository with iGit interface;
	repo, err := wc.iGit.NewRepo(name)
	if err != nil {
		return nil, err
	}

	//Set 'apps' as collaborator to created repository;
	_, err = wc.iGit.AppsAsCollaboratorTo(repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (wc *WorkflowController) CreateFolder(name string) (*d.Folder, error) {

	//Create ownCloud folder w/ iCloud interface;
	cloud, err := wc.iCloud.CreateFolder(name)
	if cloud == nil {
		return nil, err
	}

	/* [ ] Experimental call:
	wc.iCloud.ShareToExternals(cloud)
	*/

	return cloud, err
}
