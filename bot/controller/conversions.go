package controller

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/imroc/req/v3"

	d "ticket-pimp/bot/domain"
)

func (wc *WorkflowController) ThrowConversions(f io.ReadCloser, appID string, token string) *d.ConversionLog {
	c := req.C().
		SetBaseURL("https://graph.facebook.com/v15.0/").
		DevMode()

	const currency = "USD"

	r := csv.NewReader(f)

	conversionLog := d.ConversionLog{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}

		advertiser := strings.Split(record[0], ";")[0]

		params := map[string]string{
			"advertiser_id":                advertiser,
			"event":                        "CUSTOM_APP_EVENTS",
			"application_tracking_enabled": "1",
			"advertiser_tracking_enabled":  "1",
			"custom_events":                `[{"_eventName":"fb_mobile_purchase"}]`,
		}

		res, _ := c.R().
			SetQueryString(token).
			SetQueryParams(params).
			Post(appID + "/activities")

		if res.Err != nil {
			conversionLog.Advertiser = append(conversionLog.Advertiser, advertiser)
		}

	}

	return &conversionLog
}
