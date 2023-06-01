package domain

import (
	"log"
	"time"

	"github.com/imroc/req/v3"
)

type cloud struct {
	baseUrl string
	client  *req.Client
}

func NewCloud(base, user, pass string) *cloud {

	client := req.C().
		SetTimeout(5*time.Second).
		SetCommonBasicAuth(user, pass).
		SetBaseURL(base)
	return &cloud{
		baseUrl: base,
		client:  client,
	}
}

type Cloud struct {
	FolderName string
	FolderPath string
}

func (c *cloud) CreateFolder(name string) (*Cloud, error) {
	const (
		HOMEPATH = "/remote.php/dav/files/naudachu/%23mobiledev/"
		PATH     = "/apps/files/?dir=/%23mobiledev/"
	)

	cloud := Cloud{
		FolderName: name,
		FolderPath: "",
	}

	resp, err := c.client.R().
		Send("MKCOL", HOMEPATH+name)

	// Check if request failed or response status is not Ok;
	if !resp.IsSuccessState() || err != nil {
		log.Print("bad status:", resp.Status)
		log.Print(resp.Dump())
	}

	if resp.StatusCode == 201 {
		cloud.FolderPath = c.baseUrl + PATH + name
	}
	return &cloud, err
}
