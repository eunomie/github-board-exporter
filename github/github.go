package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

func (g *Github) getBytes(url string) ([]byte, error) {
	resp, err := g.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not GetString for url %s: %v", url, err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body for url %s: %v", url, err)
	}
	return bytes, nil
}

// GetString get a github path and return the body as string
func (g *Github) GetString(url string) (string, error) {
	bytes, err := g.getBytes(url)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// GetJSON fill an interface with the body as json
func (g *Github) GetJSON(url string, v interface{}) error {
	bytes, err := g.getBytes(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, v); err != nil {
		return fmt.Errorf("could not parse json for url %s: %v", url, err)
	}
	return nil
}

// Patch github url
func (g *Github) Patch(url, content string) (*http.Response, error) {
	return g.Do(http.MethodPatch, url, strings.NewReader(content))
}

func (g *Github) patchBytes(url, content string) ([]byte, error) {
	resp, err := g.Patch(url, content)
	if err != nil {
		return nil, fmt.Errorf("could not Patch for url %s: %v", url, err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body for url %s: %v", url, err)
	}
	return bytes, nil
}

// PatchJSON fill an interface with the body as json
func (g *Github) PatchJSON(url, content string, v interface{}) error {
	bytes, err := g.patchBytes(url, content)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, v); err != nil {
		return fmt.Errorf("could not parse json for url %s: %v", url, err)
	}
	return nil
}
