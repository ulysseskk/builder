package connector

import (
	"context"

	"github.com/google/go-github/v40/github"
	"github.com/ulysseskk/builder/app/common/config"
	"golang.org/x/oauth2"
)

var (
	githubConnector *github.Client
)

func initGithubConnector(conf *config.GithubConfig) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubConnector = github.NewClient(tc)
}

func GetGithubClient() *github.Client {
	return githubConnector
}
