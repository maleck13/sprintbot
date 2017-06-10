package github

import (
	"context"

	"github.com/pkg/errors"

	"strings"

	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Client struct {
	gclient *github.Client
	context context.Context
}

func NewClient(token string) *Client {
	if "" == token {
		panic("github client needs a token ")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &Client{
		context: ctx,
		gclient: github.NewClient(tc),
	}
}

func (c *Client) Repo(pURL string) string {
	parts := strings.Split(pURL, "/")
	//2,3,4
	if len(parts) < 7 {
		return ""
	}
	return parts[4]
}

func (c *Client) PRReviewed(prURL string) (bool, error) {
	//https://github.com/fheng/fh-supercore/pull/984
	parts := strings.Split(prURL, "/")
	//2,3,4
	if len(parts) < 7 {
		return false, errors.New("github pr url was badly formatted: " + prURL)
	}
	owner := parts[3]
	repo := parts[4]
	prID := parts[6]
	id, err := strconv.Atoi(prID)
	if err != nil {
		return false, errors.New("failed to parse the PR ID " + prID)
	}
	pr, _, err := c.gclient.PullRequests.Get(c.context, owner, repo, id)
	if err != nil {
		return false, errors.Wrap(err, "Checking PR unexpected error from github call ")
	}
	if pr.GetMerged() {
		return true, nil
	}
	t := pr.ClosedAt
	if t != nil {
		return true, nil
	}
	opt := &github.ListOptions{Page: 1}
	rvs, _, err := c.gclient.PullRequests.ListReviews(c.context, owner, repo, id, opt)
	if err != nil {
		return false, errors.Wrap(err, "Checking PRReviewed unexpected error from github call ")
	}
	return len(rvs) > 0, nil
}
