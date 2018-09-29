package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Result extend block
type Result struct {
	Block ExBlock `json:"result"`
}

// ExBlock extend block
type ExBlock struct {
	Confirmations int `json:"confirmations"`
}

func exGetBlock(url string, hash string) (*ExBlock, error) {
	//url := "http://obs:obs@localhost:18332"
	const pl = "{\"jsonrpc\": \"1.0\",\"id\":\"gj.com\", \"method\": \"getblock\", \"params\": [\"%s\"]}"
	str := fmt.Sprintf(pl, hash)
	payload := strings.NewReader(str)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "6cdea4cd-e1ef-20ca-58a4-2488ab6c6744")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := Result{}
	if err := json.Unmarshal([]byte(string(body)), &result); err != nil {
		return nil, err
	}
	return &result.Block, nil
}
