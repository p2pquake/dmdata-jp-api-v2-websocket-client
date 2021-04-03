package dmdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type V2Client struct {
	ApiKey string
}

type StartSocketRequest struct {
	Classifications []string `json:"classifications,omitempty"`
	Types           []string `json:"types,omitempty"`
	AppName         string   `json:"appName,omitempty"`
}

type Response struct {
	Status string `json:"status"`
}

type ListSocketResponse struct {
	Items []ListSocketItem `json:"items"`
}

type ListSocketItem struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

type StartSocketResponse struct {
	WebSocket StartSocketWebSocket `json:"websocket"`
}

type StartSocketWebSocket struct {
	URL string `json:"url"`
}

func (c *V2Client) ListSocket(status string) (*ListSocketResponse, error) {
	body, err := c.get("https://api.dmdata.jp/v2/socket?status=" + status)
	if err != nil {
		return nil, err
	}

	r := ListSocketResponse{}
	if err := c.parse(body, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

func (c *V2Client) StartSocket(classifications []string, types []string, appName string) (*StartSocketResponse, error) {
	req := StartSocketRequest{
		Classifications: classifications,
		Types:           types,
		AppName:         appName,
	}

	body, err := c.post("https://api.dmdata.jp/v2/socket", req)
	if err != nil {
		return nil, err
	}

	r := StartSocketResponse{}
	if err := c.parse(body, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

func (c *V2Client) CloseSocket(id int) error {
	body, err := c.delete("https://api.dmdata.jp/v2/socket/" + strconv.Itoa(id))
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal(body, &r); err != nil {
		return err
	}

	if r.Status == "error" {
		return fmt.Errorf("Response status error: %s", body)
	}

	return nil
}

func (c *V2Client) get(url string) ([]byte, error) {
	log.Printf("GET %s", url)

	res, err := http.Get(c.appendApiKey(url))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return body, err
}

func (c *V2Client) post(url string, v interface{}) ([]byte, error) {
	log.Printf("POST %s with %#v", url, v)

	json, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(c.appendApiKey(url), "application/json", bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return body, err
}

func (c *V2Client) delete(url string) ([]byte, error) {
	log.Printf("DELETE %s", url)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, c.appendApiKey(url), nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return body, err
}

func (c *V2Client) parse(body []byte, v interface{}) error {
	r := Response{}
	if err := json.Unmarshal(body, &r); err != nil {
		return err
	}
	if r.Status != "ok" {
		return fmt.Errorf("Response status error: %s", body)
	}

	if err := json.Unmarshal(body, &v); err != nil {
		return err
	}

	return nil
}

func (c *V2Client) appendApiKey(url string) string {
	if strings.Contains(url, "?") {
		return url + "&key=" + c.ApiKey
	}
	return url + "?key=" + c.ApiKey
}
