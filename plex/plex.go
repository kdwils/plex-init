package plex

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	http   *http.Client
	Scheme string
	Host   string
}

type token struct {
	Token string `json:"token"`
}

func NewClient(client *http.Client) *Client {
	return &Client{
		Scheme: "https",
		Host:   "plex.tv",
		http:   client,
	}
}

func (c *Client) GetServerClaimToken(ctx context.Context, authToken string) (string, error) {
	u := &url.URL{
		Scheme: c.Scheme,
		Host:   c.Host,
		Path:   "/api/claim/token.json",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error getting server claim token: could not create request: %v", err)
	}

	req.Header = getBaseHeaders(authToken)

	r, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting server claim token: error during request: %v", err)
	}

	defer r.Body.Close()

	if r.StatusCode > http.StatusOK {
		return "", fmt.Errorf("error getting server claim token: unexpected status code: %v", r.StatusCode)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("error getting server claim token: could not read response body: %v", err)
	}

	var token token
	err = json.Unmarshal(b, &token)
	if err != nil {
		return "", fmt.Errorf("error getting server claim token: could not parse response: %v\n%s", err, b)
	}

	return token.Token, nil
}

func getBaseHeaders(token string) http.Header {
	headers := make(http.Header)
	headers.Set("X-Plex-Token", token)
	return headers
}
