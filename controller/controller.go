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

type IWorkflowController interface {
	Workflow(name string) (string, error)
	CreateRepo(name string, param uint) (string, error)
	CreateFolder(name string) (*d.Cloud, error)
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
			git, gitBuild string
			cloud         *d.Cloud
		)

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			git, _ = wc.CreateRepo(issue.Key, 0)
		}()

		go func() {
			defer wg.Done()
			gitBuild, _ = wc.CreateRepo(issue.Key+"-build", 1)
		}()

		go func(ref **d.Cloud) {
			defer wg.Done()
			*ref, _ = wc.CreateFolder(issue.Key + " - " + issue.Summary)
		}(&cloud)

		wg.Wait()

		yt.UpdateIssue(issue, cloud.FolderURL, git, gitBuild)
	}
	return issue.Key, nil
}

func (wc *WorkflowController) CreateRepo(name string, param uint) (string, error) {
	//Create git repository with iGit interface;
	repo, err := wc.iGit.NewRepo(name)

	//Set 'apps' as collaborator to created repository;
	wc.iGit.AppsAsCollaboratorTo(repo)

	// Result string formatting:
	if repo != nil {
		switch param {
		case 0:
			return repo.HtmlUrl, err
		case 1:
			return fmt.Sprintf("ssh://%s/%s.git", repo.SshUrl, repo.FullName), err
		default:
			return repo.CloneUrl, err
		}
	}

	return "", err
}

func (wc *WorkflowController) CreateFolder(name string) (*d.Cloud, error) {

	//Create ownCloud folder w/ iCloud interface;
	cloud, err := wc.iCloud.CreateFolder(name)
	if cloud == nil {
		return cloud, err
	}

	/* [ ] Experimental call:
	wc.iCloud.ShareToExternals(cloud)
	*/

	return cloud, err
}
