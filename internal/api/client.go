package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultEndPoint = "https://openapi.guance.com"

var Error404 = errors.New("resource_not_found")

var EndPoints = map[string]string{
	"hangzhou":  "https://openapi.guance.com",
	"ningxia":   "https://aws-openapi.guance.com",
	"guangzhou": "https://cn4-openapi.guance.com",
	"vnet":      "https://cn5-openapi.guance.com",
	"hongkong":  "https://cn6-openapi.guance.one",
	"oregon":    "https://us1-openapi.guance.com",
	"frankfurt": "https://eu1-openapi.guance.one",
	"singapore": "https://ap1-openapi.guance.one",
}

const (
	apiPrefix = "/api/v1"

	ResourceCreate = "create"
	ResourceRead   = "read"
	ResourceUpdate = "update"
	ResourceDelete = "delete"
)

var apiURLs = map[string]map[string]string{}

type Client struct {
	Region     string
	EndPoint   string
	Token      string
	HTTPClient *http.Client
}

type Response struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
	Success   bool   `json:"success"`
	TraceID   string `json:"traceId"`
	Content   any    `json:"content"`
}

func NewClient(region, token, endPoint string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		EndPoint:   DefaultEndPoint,
		Token:      token,
	}

	if ep, ok := EndPoints[region]; ok {
		c.EndPoint = ep
	}

	if len(endPoint) != 0 {
		c.EndPoint = endPoint
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("DF-API-KEY", c.Token)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 && res.StatusCode != 404 {
		return nil, fmt.Errorf("http status code: %d, body: %s, path: %s", res.StatusCode, body, req.URL.Path)
	}

	return body, err
}

func (c *Client) post(path string, body any, content any) error {
	var (
		rb         []byte
		err        error
		bodyReader io.Reader
	)

	if body != nil {
		rb, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	if rb != nil {
		bodyReader = bytes.NewReader(rb)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", c.EndPoint, apiPrefix, path), bodyReader)
	if err != nil {
		return err
	}

	respbody, err := c.doRequest(req)
	if err != nil {
		return err
	}

	resp := Response{
		Content: content,
	}

	err = json.Unmarshal(respbody, &resp)
	if err != nil {
		return err
	}

	if !resp.Success {
		return errors.New(resp.Message)
	}

	return nil
}

func (c *Client) get(path string, content any) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.EndPoint, apiPrefix, path), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	res := Response{Content: content}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	if res.Code == 404 {
		return Error404
	}

	if !res.Success {
		return errors.New(res.Message)
	}

	return nil
}

func (c *Client) Create(resource string, item any, resp any) error {
	path, err := getResourcePath(resource, ResourceCreate)
	if err != nil {
		return fmt.Errorf("api path for create not found: %w", err)
	}
	return c.post(path, item, resp)
}

func (c *Client) Delete(resource string, key string) error {
	path, err := getResourcePath(resource, ResourceDelete)
	if err != nil {
		return fmt.Errorf("api path for delete not found: %w", err)
	}
	path = fmt.Sprintf(path, key)
	err = c.get(path, nil)

	// just ignore 404 error
	if err == Error404 {
		return nil
	}
	return err
}

func (c *Client) Read(resource string, key string, resp any) error {
	path, err := getResourcePath(resource, ResourceRead)
	if err != nil {
		return fmt.Errorf("api path for read not found: %w", err)
	}
	path = fmt.Sprintf(path, key)
	return c.get(path, resp)
}

func (c *Client) Update(resource, key string, item any, resp any) error {
	path, err := getResourcePath(resource, ResourceUpdate)
	if err != nil {
		return fmt.Errorf("api path for update not found: %w", err)
	}
	path = fmt.Sprintf(path, key)
	return c.post(path, item, resp)
}

// getResourcePath returns the API path corresponding to the specified resource and operation.
func getResourcePath(resource string, operation string) (string, error) {
	// Retrieve the URL mappings for the given resource.
	urls, ok := apiURLs[resource]
	if !ok {
		return "", errors.New("unknown resource")
	}

	path, ok := urls[operation]
	if !ok {
		return "", errors.New("create path not found")
	}

	return path, nil
}
