package asana

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	baseUrl             = "https://app.asana.com/api/1.0"
	envAsanaAccessToken = "ASANA_PERSONAL_ACCESS_TOKEN"
	envAsanaTaskListID  = "ASANA_TASK_LIST_ID"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

type Task struct {
	id        string `json:"gid,omitempty"`
	completed bool   `json:"completed,omitempty"`
}

type NextPage struct {
	offset string `json:"offset"`
	path   string `json:"path"`
	uri    string `json:"uri"`
}

type Response struct {
	data     json.RawMessage `json:"data"`
	nextPage *NextPage       `json:"next_page"`
	err      error
}

func NewClient() (*Client, error) {
	ctx := context.Background()
	accessToken := os.Getenv(envAsanaAccessToken)
	if accessToken == "" {
		return nil, fmt.Errorf("%v was not set in your environment", envAsanaAccessToken)
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})
	client := oauth2.NewClient(ctx, tokenSource)
	u, _ := url.Parse(baseUrl)
	return &Client{
		baseURL:    u,
		httpClient: client,
	}, nil
}

func (c *Client) ListAllTasks() ([]*Task, *NextPage, error) {
	taskListID := os.Getenv(envAsanaTaskListID)
	if taskListID == "" {
		return nil, nil, fmt.Errorf("%v was not set in your environment", envAsanaTaskListID)
	}
	var result []*Task
	nextPage, err := c.getAllTasks(fmt.Sprintf("/projects/%s/tasks?opt_fields=completed", taskListID), &result)
	return result, nextPage, err
}

func (c *Client) getAllTasks(path string, result *[]*Task) (*NextPage, error) {
	requestID := xid.New()
	request, err := http.NewRequest(http.MethodGet, c.getURL(path), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "%s Request error", requestID)
	}
	res, err := c.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "%s GET error", requestID)
	}
	resultData, err := c.parseResponse(res, result, requestID)
	if err != nil {
		return nil, err
	}

	return resultData.nextPage, nil
}

func (c *Client) getURL(path string) string {
	return c.baseURL.String() + path
}

func (c *Client) parseResponse(res *http.Response, result interface{}, requestID xid.ID) (*Response, error) {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	value := &Response{}
	if err := json.Unmarshal(body, value); err != nil {
		return nil, errors.Wrapf(err, "%s Unable to parse response body", requestID)
	}
	if value.data == nil {
		return nil, errors.Errorf("%s Missing data from response", requestID)
	}
	return value, c.parseResponseData(value.data, result, requestID)
}

func (c *Client) parseResponseData(data []byte, result interface{}, requestID xid.ID) error {
	if result == nil {
		return nil
	}
	if err := json.Unmarshal(data, result); err != nil {
		return errors.Wrapf(err, "%s Unable to parse response data", requestID)
	}
	return nil
}
