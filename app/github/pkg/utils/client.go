package utils

import (
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/google/go-github/v71/github"
	"io"
	"net/http"
)

type GithubApi struct {
	Default *github.Client
	Org     string

	baseUrl string
	token   string
}

func NewGithubApiClient(c *config.Config) *GithubApi {
	return &GithubApi{
		Default: github.NewClient(nil).WithAuthToken(c.Env.GetString("pat")),
		Org:     c.Env.GetString("org"),
		baseUrl: c.Env.GetString("api"),
		token:   c.Env.GetString("pat"),
	}
}

func (client *GithubApi) Request(method string, path string) ([]byte, error) {
	req, err := http.NewRequest(method, client.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+client.token)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	res, err := client.Default.Client().Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
