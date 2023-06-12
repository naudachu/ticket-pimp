package ext

import (
	"os"
	"time"
)

func NewCloud(base, user, pass string) *Client {

	client := NewClient().
		SetTimeout(5*time.Second).
		SetCommonBasicAuth(user, pass).
		SetBaseURL(base)

	return &Client{
		client,
	}
}

type Cloud struct {
	FolderName string
	FolderPath string
}

func (c *Client) CreateFolder(name string) (*Cloud, error) {

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
