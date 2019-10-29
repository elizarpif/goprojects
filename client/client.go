// client.go
package main

import (
	"bytes"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	browser  http.Client
	basePath url.URL
}

// NewClient in a constructor for *Client ?
func NewClient(client *http.Client, basePath string) (*Client, error) {
	u, err := url.Parse(basePath)
	if err != nil {
		return nil, err
	}
	return &Client{
		browser:  *client,
		basePath: *u,
	}, nil
}

type DadataRequest struct {
	Query      string `json:"query"`
	BranchType string `json:"branch_type"`
}

type DadataResponse struct {
	Suggestions []Suggestions `json: "suggestions"`
}
type Suggestions struct {
	Data Data `json:"data"`
}

type Data struct {
	Inn     string  `json:"inn"`
	Kpp     string  `json:"kpp"`
	Orgn    string  `json:"orgn"`
	Address Address `json:"address"`
	Name    Name    `json:"short"`
	Opf     Name    `json:"short"`
}
type Name struct { //МТС
	Short string `json:"name"`
}
type Opf struct { //ПАО
	Short string `json:"opf"`
}

type Address struct {
	Region         string `json: "region"`
	City           string `json: "city"`
	PostalCode     string `json: "postal_code"`
	StreetWithType string `json: "street_with_type"`
	House          string `json: "house"`
	Flat           string `json: "flat"`
}

func (c *Client) SendReduest(dadataRequest *DadataRequest) (*DadataResponse, error) {
	//	curl -X POST   -H "Content-Type: application/json"   -H "Accept: application/json"
	//	-H "Authorization: Token 5eff14ad466a316f4e23bd4dff7f631880c44eed"
	//	-d '{ "query": "7740000076", "branch_type":"MAIN"}'
	//  https://suggestions.dadata.ru/suggestions/api/4_1/rs/findById/party | jq .
	//logger := logging.GetLogger(ctx)
	dadataURL := fmt.Sprintf("https://suggestions.dadata.ru/suggestions/api/4_1/rs/findById/party ")

	body, err := json.Marshal(dadataRequest)
	if err != nil {
		//logger.WithError(err).Error("create request failed")
		return nil, err
	}
	bodyReader := bytes.NewReader(body)
	request, err := http.NewRequest(http.MethodPost, dadataURL, bodyReader)
	if err != nil {
		fmt.Println(err, "err with newRequest")
		return nil, err
	}
	request.Header.Add("Authorization", "Token 5eff14ad466a316f4e23bd4dff7f631880c44eed")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	res, err := c.browser.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			fmt.Println("can not close response body")
			//logger.WithError(err).Error("can not close response body")
		}
	}()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//logger.WithError(err).Error("read dadata response failed")
		return nil, err
	}

	dadataResponse := new(DadataResponse)
	err = json.Unmarshal(responseBody, dadataResponse)
	if err != nil {
		//logger.WithError(err).Error("unmarshal dadata response failed")
		return nil, err
	}
	return dadataResponse, nil

}

func main() {
	dadataClient, err := NewClient(&http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: time.Duration(10) * time.Second,
	}, "http://localhost:3000")
	if err != nil {
		//logger.WithError(err).Error("create dadata client failed")
		fmt.Println("create dadata client failed")
		fmt.Println(err)
	}

	dadatareq := &DadataRequest{
		Query:      "7740000076",
		BranchType: "MAIN",
	}

	dadataResponce, err := dadataClient.SendReduest(dadatareq)
	if err != nil {
		fmt.Println("error dataresponce")
	}
	fmt.Println(dadataResponce.Suggestions[0].Data)
}
