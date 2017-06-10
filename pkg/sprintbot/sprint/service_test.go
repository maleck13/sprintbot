package sprint_test

import (
	"strings"
	"testing"

	"fmt"

	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/maleck13/sprintbot/pkg/sprintbot/sprint"
)

type mockIssueFinder struct {
	issues []mockIssue
	err    error
}

func (mf *mockIssueFinder) FindUnresolvedOnBoard(boardName, sprint string) (*sprintbot.IssueList, error) {
	var ret = make([]sprintbot.IssueState, len(mf.issues))
	for i, d := range mf.issues {
		ret[i] = d
	}
	return sprintbot.NewIssueList(ret), mf.err
}

type mockRepoChecker struct {
	issues []mockIssue
}

func (mrc *mockRepoChecker) PRReviewed(prURL string) (bool, error) {
	for _, i := range mrc.issues {
		for _, pr := range i.PRS() {
			if pr == prURL {
				return i.prReviewed, nil
			}
		}
	}
	return false, nil
}
func (mrc *mockRepoChecker) Repo(pURL string) string {
	parts := strings.Split(pURL, "/")
	fmt.Println(len(parts))
	//2,3,4
	if len(parts) < 7 {
		return ""
	}
	return parts[4]
}

type mockIssue struct {
	pr         string
	link       string
	prReviewed bool
	state      string
}

func (mi mockIssue) PRS() []string {
	return []string{mi.pr}
}
func (mi mockIssue) Link() string {
	return mi.link
}
func (mi mockIssue) State() string {
	return mi.state
}

func TestSprintNext(t *testing.T) {
	cases := []struct {
		Name        string
		ExpectError bool
		NumIssues   int
		Issues      []mockIssue
		Err         error
	}{
		{
			Name:        "should find and return no issues",
			ExpectError: false,
			NumIssues:   0,
			Issues:      []mockIssue{},
			Err:         nil,
		},
		{
			Name:        "should find and return pr issues",
			ExpectError: false,
			NumIssues:   1,
			Issues: []mockIssue{
				mockIssue{pr: "github.com/some/pr/20", link: "http://jira.com", prReviewed: false, state: sprintbot.IssueStateReadyForQA},
			},
			Err: nil,
		},
		{
			Name:        "should not find and return pr issues for ignored repo",
			ExpectError: false,
			NumIssues:   0,
			Issues: []mockIssue{
				mockIssue{pr: "https://gitlab.cee.redhat.com/red-hat-mobile-application-platform-documentation/RHMAPDocsNG/merge_requests/181", link: "http://jira.com", prReviewed: false, state: sprintbot.IssueStateClosed},
			},
			Err: nil,
		},
		{
			Name:        "should find and return ready for qe  issues",
			ExpectError: false,
			NumIssues:   1,
			Issues: []mockIssue{
				mockIssue{pr: "github.com/some/pr/20", link: "http://jira.com", prReviewed: true, state: sprintbot.IssueStateReadyForQA},
			},
			Err: nil,
		},
		{
			Name:        "should find and return Open  issues",
			ExpectError: false,
			NumIssues:   1,
			Issues: []mockIssue{
				mockIssue{pr: "github.com/some/pr/20", link: "http://jira.com", prReviewed: true, state: sprintbot.IssueStateOpen},
			},
			Err: nil,
		},
	}

	for _, tc := range cases {
		isf := &mockIssueFinder{
			issues: tc.Issues,
			err:    tc.Err,
		}
		rc := &mockRepoChecker{issues: tc.Issues}
		s := sprintbot.Sprint{Name: "testSprint", Board: "testBoard"}
		service := sprint.NewService(isf, rc, &s)
		service.IgnoredRepos = []string{"RHMAPDocsNG", "fhcap", "fh-openshift-templates", "fh-core-openshift-templates"}
		t.Run(tc.Name, func(t *testing.T) {
			ni, err := service.Next()
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none ")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got one  %s ", err)
			}
			if ni == nil {
				t.Fatal("Next issues was returned nil ")
			}
			if len(ni.Issues) != tc.NumIssues {
				t.Fatalf("expected %v issues but got  %v ", tc.NumIssues, len(ni.Issues))
			}
		})
	}

}
