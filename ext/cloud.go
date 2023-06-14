package ext

import (
	"os"
	"ticket-pimp/helpers"
	"time"
)

type Cloud struct {
	//[ ] out in separate domain struct
	FolderName string
	FolderPath string
	*Client
}

type ICloud interface {
	CreateFolder(name string) (*Cloud, error)
}

func NewCloud(base, user, pass string) *Cloud {

	client := NewClient().
		SetTimeout(5*time.Second).
		SetCommonBasicAuth(user, pass).
		SetBaseURL(base)

	return &Cloud{
		FolderName: "",
		FolderPath: "",
		Client: &Client{
			client,
		},
	}
}

func (c *Cloud) CreateFolder(name string) (*Cloud, error) {

	name = helpers.GitNaming(name)

	cloud := Cloud{
		FolderName: name,
		FolderPath: "",
	}

	pathName := os.Getenv("HOMEPATH") + name

	resp, err := c.R().
		Send("MKCOL", pathName)

	if resp.IsSuccessState() {
		cloud.FolderPath = c.BaseURL + os.Getenv("FOLDER_PATH") + name

	}

	return &cloud, err
}
