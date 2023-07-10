package controller

import (
	"ticket-pimp/internal/domain"
	"ticket-pimp/internal/extapi"
)

type Cloud struct {
	api extapi.ICloud
}

func NewCloudController(url, user, pass string) *Cloud {
	return &Cloud{
		extapi.NewCloudClient(url, user, pass),
	}
}

type CloudCreator interface {
	CreateFolder(name string) (*domain.Folder, error)
}

func (c *Cloud) CreateFolder(name string) (*domain.Folder, error) {

	//Create ownCloud folder w/ iCloud interface;
	cloud, err := c.api.CreateFolder(name)
	if cloud == nil {
		return nil, err
	}

	return cloud, err
}
