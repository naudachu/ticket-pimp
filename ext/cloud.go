package ext

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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

type MultistatusObj struct {
	XMLName     xml.Name `xml:"multistatus"`
	Multistatus struct {
		XMLName  xml.Name `xml:"response"`
		Propstat struct {
			XMLName xml.Name `xml:"propstat"`
			Prop    struct {
				XMLName xml.Name `xml:"prop"`
				FileID  struct {
					XMLName xml.Name `xml:"fileid"`
					ID      string   `xml:",chardata"`
				}
			}
		}
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

		xmlFile, err := ioutil.ReadFile("./fileid.xml")

		if err != nil {
			fmt.Println(err)
			return nil, err // fix this return;
		}

		resp, _ := c.R().
			SetBody(xmlFile).
			Send("PROPFIND", os.Getenv("HOMEPATH")+os.Getenv("ROOTDIR")+cloud.FolderName)

		id, err := getFileIDFromRespBody(resp.Bytes())

		if err != nil {
			log.Print(err) // [ ] Если тут проблема - надо пытаться засетать полную ссылку
		}

		cloud.PrivateURL = os.Getenv("CLOUD_BASE_URL") + "/f/" + strconv.Itoa(id)
	}

	return &cloud, err
}

func getFileIDFromRespBody(str []byte) (int, error) {

	var multi MultistatusObj

	err := xml.Unmarshal(str, &multi)
	if err != nil {
		return 0, fmt.Errorf("XML Unmarshal error: %v", err)
	}

	id, err := strconv.Atoi(multi.Multistatus.Propstat.Prop.FileID.ID)
	if err != nil {
		return 0, fmt.Errorf("FileID str to int convertion error: %v", err)
	}

	return id, nil
}

func (c *Cloud) ShareToExternals(cloud *domain.Cloud) (*domain.Cloud, error) {
	return nil, nil
}
