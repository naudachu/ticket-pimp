package ext

import (
	"fmt"
	"os"
	"strconv"
	d "ticket-pimp/domain"
	"ticket-pimp/helpers"
	"time"
)

type Cloud struct {
	*Client
}

type ICloud interface {
	CreateFolder(name string) (*d.Folder, error)
	ShareToExternals(cloud *d.Folder) (*d.Folder, error)
}

func NewCloud(base, user, pass string) *Cloud {

	client := NewClient().
		SetTimeout(5*time.Second).
		SetCommonBasicAuth(user, pass).
		SetBaseURL(base)

	return &Cloud{
		Client: &Client{
			client,
		},
	}
}

func (c *Cloud) CreateFolder(name string) (*d.Folder, error) {
	rootDir := os.Getenv("ROOTDIR")
	user := os.Getenv("CLOUD_USER")

	davPath := "/remote.php/dav/files/"
	parentPath := "/apps/files/?dir="

	name = helpers.GitNaming(name)

	cloud := d.Folder{
		Title:      name,
		PrivateURL: "",
	}

	requestPath := davPath + user + rootDir + name

	cloud.PathTo = parentPath + rootDir + name

	resp, _ := c.R().
		Send("MKCOL", requestPath)

	if resp.IsSuccessState() {
		// Set stupid URL to the d entity
		cloud.PrivateURL = c.BaseURL + cloud.PathTo

		// Try to set short URL to the d entity
		if err := c.setPrivateURL(requestPath, &cloud); err != nil {
			return &cloud, err
		}
	}

	return &cloud, nil
}

func (c *Cloud) setPrivateURL(requestPath string, cloud *d.Folder) error {

	payload := []byte(`<?xml version="1.0"?><a:propfind xmlns:a="DAV:" xmlns:oc="http://owncloud.org/ns"><a:prop><oc:fileid/></a:prop></a:propfind>`)

	// Deprecated: Read XML file
	/*
		xmlFile, err := ioutil.ReadFile("./fileid.xml") // moved into this method as a string..

		if err != nil {
			return fmt.Errorf("request xml file error: %v", err)
		}
	*/

	resp, _ := c.R().
		SetBody(payload).
		Send("PROPFIND", requestPath)

	if resp.Err != nil {
		return resp.Err
	}

	id := helpers.GetFileIDFromRespBody(resp.Bytes())

	if id == 0 {
		return fmt.Errorf("unable to get fileid")
	}

	cloud.PrivateURL = c.BaseURL + "/f/" + strconv.Itoa(id)

	return nil
}

func (c *Cloud) ShareToExternals(cloud *d.Folder) (*d.Folder, error) {
	return nil, nil
}
