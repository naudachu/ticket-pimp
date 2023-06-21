package ext

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"ticket-pimp/domain"
	"ticket-pimp/helpers"
	"time"
)

type Cloud struct {
	*Client
}

type ICloud interface {
	CreateFolder(name string) (*domain.Cloud, error)
	ShareToExternals(cloud *domain.Cloud) (*domain.Cloud, error)
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

func (c *Cloud) CreateFolder(name string) (*domain.Cloud, error) {
	rootDir := os.Getenv("ROOTDIR")

	name = helpers.GitNaming(name)

	cloud := domain.Cloud{
		FolderName: name,
		FolderURL:  "",
	}

	requestPath := os.Getenv("HOMEPATH") + rootDir + name
	cloud.FolderPath = os.Getenv("FOLDER_PATH") + rootDir + name

	resp, err := c.R().
		Send("MKCOL", requestPath)

	if resp.IsSuccessState() {

		cloud.FolderURL = c.BaseURL + cloud.FolderPath

		/*
			type ResponseObj struct {
				Multistatus struct {
					Response struct {
						Href struct {
							Propstat struct {
								Prop struct {
									FileID int `json:"oc:fileid"`
								} `json:"d:prop"`
							} `json:"d:propstat"`
						} `json:"d:href"`
					} `json:"d:response"`
				} `json:"d:multistatus"`
			}*/

		type ResponseObj struct {
			XMLName     xml.Name `xml:"d:multistatus"`
			Multistatus struct {
				XMLName  xml.Name `xml:"d:multistatus"`
				Response struct {
					Href struct {
						Propstat struct {
							Prop struct {
								FileID string `xml:"oc:fileid"`
							} `xml:"d:prop"`
						} `xml:"d:propstat"`
					} `xml:"d:href"`
				} `xml:"d:response"`
			} `xml:"d:multistatus"`
		}

		xmlFile, err := ioutil.ReadFile("./fileid.xml")

		if err != nil {
			fmt.Println(err)
			return nil, err // fix this return;
		}

		var id ResponseObj

		resp, _ := c.R().
			SetBody(xmlFile).
			Send("PROPFIND", os.Getenv("HOMEPATH")+os.Getenv("ROOTDIR")+cloud.FolderName)

		xmlEncodingErr := resp.UnmarshalXml(&id)
		if xmlEncodingErr != nil {
			log.Print(err)
		}

		log.Print(resp)

	}

	return &cloud, err
}

func (c *Cloud) ShareToExternals(cloud *domain.Cloud) (*domain.Cloud, error) {
	return nil, nil
}
