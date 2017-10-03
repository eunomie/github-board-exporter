package github

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Github contains all logic to call Github's api
type Github struct {
	token  string
	client *http.Client
}

// NewGithub creates a Github handler
func NewGithub(token string) (*Github, error) {
	client := http.DefaultClient
	return &Github{
		token:  token,
		client: client,
	}, nil
}

func (g *Github) authURL(githubURL string) (*url.URL, error) {
	u, err := url.Parse(githubURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse github url for url %s: %v", githubURL, err)
	}
	query := u.Query()
	query.Set("access_token", g.token)
	u.RawQuery = query.Encode()
	return u, nil
}

func (g *Github) request(method, url string, body io.Reader) (*http.Request, error) {
	u, err := g.authURL(url)
	if err != nil {
		return nil, fmt.Errorf("could not create url for path %s: %v", url, err)
	}
	req, err := http.NewRequest(method, u.String(), body)
	req.Header.Add("Accept", "application/vnd.github.inertia-preview+json")
	return req, nil
}

func (g *Github) do(req *http.Request) (*http.Response, error) {
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute request %v: %v", req, err)
	}
	return resp, nil
}

// Do a request to github
func (g *Github) Do(method, url string, body io.Reader) (*http.Response, error) {
	req, err := g.request(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("could not get %s request for url %s: %v", method, url, err)
	}
	return g.do(req)
}

// Get a github path
func (g *Github) Get(url string) (*http.Response, error) {
	return g.Do(http.MethodGet, url, http.NoBody)
}

// GetString get a github path and return the body as string
func (g *Github) GetString(url string) (string, error) {
	resp, err := g.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not GetString for url %s: %v", url, err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read body for url %s: %v", url, err)
	}
	return string(bytes), nil
}
