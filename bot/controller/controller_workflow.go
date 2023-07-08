package controller

import (
	"fmt"
	"sync"
	d "ticket-pimp/bot/domain"
)

func (wc *WorkflowController) Workflow(name string) (string, error) {
	yt := wc.iYouTrack

	projectID, err := yt.GetProjectIDByName("APP")
	if err != nil {
		return "", err
	}

	// Create an issue at the available project with the provided name
	issue, err := yt.CreateIssue(projectID, name, "")

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
