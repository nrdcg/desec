package desec

import (
	"net/http"
	"net/url"
	"path"
)

const defaultBaseURL = "https://desec.io/api/v1/"

type service struct {
	client *Client
	token  string
}

// Client deSEC API client.
type Client struct {
	// HTTP client used to communicate with the API.
	HTTPClient *http.Client

	// Base URL for API requests.
	BaseURL string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the deSEC API.
	Tokens  *TokensService
	Records *RecordsService
	Domains *DomainsService
}

// NewClient creates a new Client.
func NewClient(token string) *Client {
	client := &Client{HTTPClient: http.DefaultClient, BaseURL: defaultBaseURL}

	client.common.client = client
	client.common.token = token

	client.Tokens = (*TokensService)(&client.common)
	client.Records = (*RecordsService)(&client.common)
	client.Domains = (*DomainsService)(&client.common)

	return client
}

func (c *Client) createEndpoint(parts ...string) (*url.URL, error) {
	return createEndpoint(c.BaseURL, parts)
}

func createEndpoint(baseURL string, parts []string) (*url.URL, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	endpoint, err := base.Parse(path.Join(base.Path, path.Join(parts...)))
	if err != nil {
		return nil, err
	}

	endpoint.Path += "/"

	return endpoint, nil
}
