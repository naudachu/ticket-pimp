package ext

import (
	"log"
	"time"
)

type Coda struct {
	*Client
}

type ICoda interface {
	ListDocs()
}

func NewCodaClient() *Coda {

	client := NewClient().
		SetTimeout(5 * time.Second).
		SetCommonBearerAuthToken("f54477f0-98ca-4285-844f-9fa2ef34475d").
		SetBaseURL("https://coda.io/apis/v1")

	return &Coda{
		Client: &Client{
			client,
		},
	}
}

func (c *Coda) ListDocs() {

	const tableID = "grid-obBN3tWdeh"
	const docID = "Ic3IZpQ3Wk"

	//var i []RespObj

	resp, _ := c.R().
		SetQueryParam("tableTypes", "table").
		//SetSuccessResult(&i).
		Get("/docs/" + docID + "/tables/" + tableID + "/rows")

	if resp.Err != nil {
		log.Print(resp.Err)
		return
	}

	log.Print(resp)

}
