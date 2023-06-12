package controller

import (
	"fmt"
	"os"
	"sync"
	"ticket-pimp/ext"
)

func Workflow(name string) (string, error) {
	yt := ext.NewYT(os.Getenv("YT_URL"), os.Getenv("YT_TOKEN"))

	projects, err := yt.GetProjects()

	if err != nil {
		return "", err
	}

	issue, err := yt.CreateIssue(projects[0].ID, name)

	if err != nil {
		return "", err
	}

	if issue != nil {
		var (
			git, gitBuild, folder string
		)

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			git, _ = CreateRepo(issue.Key, 0)
		}()

		go func() {
			defer wg.Done()
			gitBuild, _ = CreateRepo(issue.Key+"-build", 1)
		}()

		go func() {
			defer wg.Done()
			folder = CreateFolder(issue.Key + " - " + issue.Summary)
		}()

		wg.Wait()

		yt.UpdateIssue(issue, folder, git, gitBuild)
	}
	return issue.Key, nil
}

func CreateRepo(name string, param uint) (string, error) {
	gb := ext.NewGit(os.Getenv("GIT_BASE_URL"), os.Getenv("GIT_TOKEN"))
	repo, err := gb.NewRepo(name)
	gb.AppsAsCollaboratorTo(repo)

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

func CreateFolder(name string) string {
	oc := ext.NewCloud(os.Getenv("CLOUD_BASE_URL"), os.Getenv("CLOUD_USER"), os.Getenv("CLOUD_PASS"))

	cloud, _ := oc.CreateFolder(name)
	if cloud != nil {
		return cloud.FolderPath
	}
	return "no-folder"
}
