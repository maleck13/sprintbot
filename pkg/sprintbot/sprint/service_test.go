package sprint_test

import (
	"testing"

	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/maleck13/sprintbot/pkg/sprintbot/sprint"
)

type mockIssueFinder struct {
	issues []mockIssue
	err    error
}

func (mf *mockIssueFinder) FindUnresolved() (*sprintbot.IssueList, error) {
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
		if i.PR() == prURL {
			return i.prReviewed, nil
		}
	}
	return false, nil
}

type mockIssue struct {
	pr         string
	link       string
	prReviewed bool
	state      string
}

func (mi mockIssue) PR() string {
	return mi.pr
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
		service := sprint.NewService(isf, rc)
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
