package ext

import (
	"time"

	"github.com/imroc/req/v3"
)

type cloud struct {
	*req.Client
}

func NewCloud(base, user, pass string) *cloud {

	client := NewClient().
		SetTimeout(5*time.Second).
		SetCommonBasicAuth(user, pass).
		SetBaseURL(base)

	return &cloud{
		client,
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

	pathName := HOMEPATH + name

	resp, err := c.R().
		Send("MKCOL", pathName)

	if resp.IsSuccessState() {
		cloud.FolderPath = c.BaseURL + PATH + name
	}

	return &cloud, err
}
