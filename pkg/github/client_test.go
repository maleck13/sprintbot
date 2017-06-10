package github_test

import (
	"testing"

	"os"

	"github.com/maleck13/sprintbot/pkg/github"
)

func TestPRReviewed(t *testing.T) {
	var (
		token = os.Getenv("GITHUB_TOKEN")
		pr    = os.Getenv("GITHUB_PR")
	)
	if token == "" || pr == "" {
		t.Skip("skipping as we need a real token and pr to run this test ")
	}
	c := github.NewClient(token)
	is, err := c.PRReviewed(pr)
	if err != nil {
		t.Fatalf("did not expect error checking pr %s", err)
	}
	if !is {
		t.Fatal("did not expect the PR to be reviewed ")
	}
}

func TestRepo(t *testing.T) {
	cases := []struct {
		Name string
		PR   string
		Repo string
	}{
		{
			Name: "test get repo name ok",
			PR:   "https://github.com/fheng/a/pull/983",
			Repo: "a",
		},
		{
			Name: "test get repo name ok",
			PR:   "https://github.com/fheng/b/pull/983",
			Repo: "b",
		},
		{
			Name: "test get repo name ok",
			PR:   "https://github.com/fheng/a/",
			Repo: "",
		},
		{
			Name: "gitlab test",
			PR:   "https://gitlab.cee.redhat.com/red-hat-mobile-application-platform-documentation/RHMAPDocsNG/merge_requests/181",
			Repo: "RHMAPDocsNG",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			c := github.NewClient("t")
			if tc.Repo != c.Repo(tc.PR) {
				t.Fatalf("expected to find repo %s but got %s", tc.Repo, c.Repo(tc.PR))
			}
		})
	}
}
