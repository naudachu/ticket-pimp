package controller

import (
	"fmt"
	"sync"
	d "ticket-pimp/bot/domain"
)

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
