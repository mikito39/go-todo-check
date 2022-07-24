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
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

type Task struct {
	ID        string `json:"gid"`
	Completed bool   `json:"completed"`
}

type nextPage struct {
	offset string `json:"offset"`
	path   string `json:"path"`
	uri    string `json:"uri"`
}

type Response struct {
	Data     []*Task   `json:"data,omitempty"`
	nextPage *nextPage `json:"next_page,omitempty"`
	Err      error
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

func (c *Client) ListAllTasks(taskListID string) (chan *Response, error) {
	resChan := make(chan *Response)
	go func() {
		defer close(resChan)
		path := fmt.Sprintf("/projects/%s/tasks?opt_fields=completed", taskListID)
		for {
			requestID := xid.New()
			res, err := c.getResponse(path, requestID)
			if err != nil {
				resChan <- &Response{Err: err}
				return
			}
			page, err := c.parseResponse(res, requestID)
			if err != nil {
				resChan <- &Response{Err: err}
				return
			}

			resChan <- page

			if np := page.nextPage; np != nil && np.path == "" {
				path = np.path + "&opt_fields=completed"
			} else {
				break
			}
		}
	}()
	return resChan, nil
}

func (c *Client) getResponse(path string, requestID xid.ID) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, c.getURL(path), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "%s Request error", requestID)
	}
	res, err := c.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "%s GET error", requestID)
	}
	return res, nil
}

func (c *Client) getURL(path string) string {
	return c.baseURL.String() + path
}

func (c *Client) parseResponse(res *http.Response, requestID xid.ID) (*Response, error) {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	page := new(Response)
	if err := json.Unmarshal(body, page); err != nil {
		page.Err = err
	}
	if page.Data == nil {
		return nil, errors.Errorf("%s Missing data from response", requestID)
	}
	return page, nil
}
