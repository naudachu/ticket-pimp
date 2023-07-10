package extapi

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type CommonClient struct {
	*req.Client
}

func NewClient() *CommonClient {
	return &CommonClient{req.C().
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil {
				if dump := resp.Dump(); dump != "" {
					resp.Err = fmt.Errorf("%s\nraw content:\n%s", resp.Err.Error(), resp.Dump())
				}
				return nil // Skip the following logic if there is an underlying error.
			}

			if !resp.IsSuccessState() {
				resp.Err = fmt.Errorf("bad response, raw content:\n%s", resp.Dump())
				return nil
			}
			return nil
		}),
	}
}
