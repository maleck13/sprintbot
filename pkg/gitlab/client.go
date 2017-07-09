package gitlab

import (
        "github.com/xanzy/go-gitlab"
	"github.com/pkg/errors"
	"context"
	"strconv"
	"golang.org/x/oauth2"
	"strings"
)


type Client struct {
	gclient *gitlab.Client
	context context.Context
}

func NewClient(token string) *Client {
	if "" == token{
		panic("gitlab client needs a token")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		gclient: gitlab.NewClient(nil, tc),
		context: ctx,
	}

}

func (c *Client) PRReviewed(prURL string) (bool, error) {
	//https://gitlab.cee.redhat.com/red-hat-mobile-application-platform-documentation/RHMAPDocsNG/merge_requests/209
	//protocol://domain/org/repo/merge_requests/number
	parts := strings.Split(prURL, "/")

	//owner := parts[3]
	repo := parts[4]
	mrID := parts[6]
	id, err := strconv.Atoi(mrID)
	if err != nil {
		return false, errors.New("failed to parse the PR ID " + mrID)
	}
	//gitlab has the concept of merge request approvals
	//https://docs.gitlab.com/ee/user/project/merge_requests/merge_request_approvals.html#configuring-approvals
	//not currently in use so using pr status e.g. "closed" for now
	mr, _, err := c.gclient.MergeRequests.GetMergeRequest(repo,id)
	if err != nil {
		return false, errors.Wrap(err, "Unexpected Error from gitlab call")
	}
	//if state not equals to closed
	if(strings.Compare(mr.State, "closed")== -1){
		return true, nil
	}
	// otherwise
	return false, nil

}
