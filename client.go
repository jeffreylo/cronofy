package cronofy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const baseURL = "https://api.cronofy.com/v1"

func (c *Client) getAPIEndpoint(path string) string {
	return baseURL + path
}

type Client struct {
	accessToken string
	baseURL     string
	client      *http.Client
}

type Config struct {
	AccessToken, BaseURL string
}

func NewClient(cfg *Config) *Client {
	apiURL := baseURL
	if cfg.BaseURL != "" {
		apiURL = cfg.BaseURL
	}
	return &Client{
		accessToken: cfg.AccessToken,
		baseURL:     apiURL,
	}
}

func (c *Client) httpClient() *http.Client {
	if c.client == nil {
		var netTransport = &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		}
		c.client = &http.Client{
			Timeout:   time.Second * 10,
			Transport: netTransport,
		}
	}
	return c.client
}

// Get executes an http.MethodGet for the specified URL, unmarshalling
// its response to result.
func (c *Client) get(ctx context.Context, reqURL string, result interface{}) error {
	return c.do(ctx, http.MethodGet, reqURL, result)
}

func (c *Client) do(ctx context.Context, method, reqURL string, result interface{}) (err error) {
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return errors.Wrap(err, "new request failed")
	}

	req = req.WithContext(ctx)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer func() {
		if rErr := resp.Body.Close(); rErr != nil && err == nil {
			err = errors.Wrap(rErr, "failed to close response body")
		}
	}()

	if resp.StatusCode >= 400 {
		return &responseError{
			url:        req.URL,
			statusCode: resp.StatusCode,
		}
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "reading response body failed")
	}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return errors.Wrapf(err, "could not unmarshal body=%s", string(bytes))
	}
	return nil
}

func (c *Client) GetCalendars() ([]*Calendar, error) {
	var res struct {
		Calendars []*Calendar `json:"calendars"`
	}
	if err := c.get(context.TODO(), c.getAPIEndpoint("/calendars"), &res); err != nil {
		return nil, err
	}
	return res.Calendars, nil
}

func (c *Client) GetEvents(options *EventsRequest) (*EventsResponse, error) {
	var res *EventsResponse
	v, _ := query.Values(options)
	if err := c.get(context.TODO(), c.getAPIEndpoint("/events?"+v.Encode()), &res); err != nil {
		return nil, err
	}
	return res, nil
}
