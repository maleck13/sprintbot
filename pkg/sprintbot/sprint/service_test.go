package sprint_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/maleck13/sprintbot/pkg/sprintbot/sprint"
)

// note these mocks should move to a separate package if there comes a time where they need to be reused.
type mockIssueEditorFinder struct {
	called map[string]int
	issues []mockIssue
	err    error
}

func newMockIssueEditorFinder(issues []mockIssue, err error) *mockIssueEditorFinder {
	return &mockIssueEditorFinder{
		called: make(map[string]int),
		issues: issues,
		err:    err,
	}
}

func (mf *mockIssueEditorFinder) IssueHost() string {
	mf.called["IssueHost"]++
	return ""
}

func (mf *mockIssueEditorFinder) IssuesForBoard(boardName, sprint string) (*sprintbot.IssueList, error) {
	mf.called["IssuesForBoard"]++
	var ret = make([]sprintbot.IssueState, len(mf.issues))
	for i, d := range mf.issues {
		ret[i] = d
	}
	return sprintbot.NewIssueList(ret), mf.err
}

func (mf *mockIssueEditorFinder) AddComment(id, c string) error {
	mf.called["AddComment"]++
	return nil
}

func (mf *mockIssueEditorFinder) Called(method string) int {
	return mf.called[method]
}

type mockSprintFinder struct {
	s      *sprintbot.JiraSprint
	err    error
	called map[string]int
}

func newMockSprintFinder(s *sprintbot.JiraSprint, err error) *mockSprintFinder {
	return &mockSprintFinder{
		s:      s,
		err:    err,
		called: make(map[string]int),
	}
}

func (md *mockSprintFinder) SprintForBoard(sprintName, boardName string) (*sprintbot.JiraSprint, error) {
	md.called["SprintForBoard"]++
	return md.s, md.err
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
	//2,3,4
	if len(parts) < 7 {
		return ""
	}
	return parts[4]
}

type mockIssue struct {
	pr          string
	link        string
	prReviewed  bool
	state       string
	storyPoints int
}

func (mi mockIssue) StoryPoints() int {
	return mi.storyPoints
}

func (mi mockIssue) RemovePR(pr string) {

}

func (mi mockIssue) PRS() []string {
	return []string{mi.pr}
}
func (mi mockIssue) Link(host string) string {
	return mi.link
}
func (mi mockIssue) State() string {
	return mi.state
}

func (mi mockIssue) Description() string {
	return ""
}

func (mi mockIssue) ID() string {
	return ""
}

type mockIssueRepo struct {
	next    *sprintbot.NextIssues
	issues  map[string][]sprintbot.IssueState
	err     error
	comment string
	called  map[string]int
}

func newMockIssueRepo(next *sprintbot.NextIssues, comment string, err error, issues map[string][]sprintbot.IssueState) *mockIssueRepo {
	return &mockIssueRepo{
		next:    next,
		issues:  issues,
		err:     err,
		comment: comment,
		called:  make(map[string]int),
	}
}

func (mr *mockIssueRepo) FindIssuesInState(state string) ([]sprintbot.IssueState, error) {
	mr.called["FindIssuesInState"]++
	return mr.issues[state], mr.err
}
func (mr *mockIssueRepo) FindIssuesNotInStates(states []string) ([]sprintbot.IssueState, error) {
	var contains = func(k string) bool {
		for _, s := range states {
			if k == s {
				return true
			}
		}
		return false
	}
	mr.called["FindIssuesNotInState"]++
	ret := []sprintbot.IssueState{}
	for k, v := range mr.issues {
		if !contains(k) {
			for _, i := range v {
				ret = append(ret, i)
			}
		}
	}
	return ret, mr.err
}

func (mr *mockIssueRepo) SaveIssuesInState(state string, issues []sprintbot.IssueState) error {
	mr.called["SaveIssuesInState"]++
	return mr.err
}

func (mr *mockIssueRepo) SaveNext(next *sprintbot.NextIssues) error {
	mr.called["SaveNext"]++
	return mr.err
}
func (mr *mockIssueRepo) FindNext() (*sprintbot.NextIssues, error) {
	mr.called["FindNext"]++
	return mr.next, mr.err
}
func (mr *mockIssueRepo) SaveCommented(id string, commentID string) error {
	mr.called["SaveCommented"]++
	return mr.err
}
func (mr *mockIssueRepo) FindCommentOnIssue(id string, commentID string) (string, error) {
	mr.called["FindCommentOnIssue"]++
	return mr.comment, mr.err
}
func (mr *mockIssueRepo) Called(method string) int {
	return mr.called[method]
}

func TestSprintNext(t *testing.T) {
	cases := []struct {
		Name        string
		Next        *sprintbot.NextIssues
		ExpectError bool
		Err         error
	}{
		{
			Name: "test find next happy",
			Next: &sprintbot.NextIssues{
				Message: "issues found",
				Issues:  []*sprintbot.Issue{},
			},
			ExpectError: false,
		},
		{
			Name:        "test find next error",
			Next:        nil,
			Err:         errors.New("failed to find next"),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		isf := &mockIssueEditorFinder{}
		sprintFinder := newMockSprintFinder(nil, nil)
		t.Run(tc.Name, func(t *testing.T) {
			rc := &mockRepoChecker{}
			s := sprintbot.Sprint{Name: "testSprint", Board: "testBoard"}
			logger := logrus.StandardLogger()
			issueRepo := newMockIssueRepo(tc.Next, "", tc.Err, nil)
			service := sprint.NewService(isf, sprintFinder, rc, issueRepo, &s, logger)
			ni, err := service.Next()
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got %s", err)
			}
			if ni == nil && tc.Next != nil {
				t.Fatalf("did not expect next issues to be nil")
			}
			if tc.Next != nil {
				if len(tc.Next.Issues) != len(ni.Issues) {
					t.Fatalf("expected %v issues but got %v", len(tc.Next.Issues), len(ni.Issues))
				}
			}

		})
	}
}

func TestSprintSync(t *testing.T) {
	cases := []struct {
		Name                   string
		ExpectError            bool
		NumIssues              int
		Issues                 []mockIssue
		IssueRepoCalls         map[string]int
		IssueEditorFinderCalls map[string]int
		Err                    error
	}{
		{
			Name:           "should find and return no issues",
			ExpectError:    false,
			NumIssues:      0,
			Issues:         []mockIssue{},
			Err:            nil,
			IssueRepoCalls: map[string]int{"SaveNext": 1},
		},
		{
			Name:        "should find and return pr issues",
			ExpectError: false,
			NumIssues:   1,
			Issues: []mockIssue{
				mockIssue{pr: "github.com/some/pr/20", link: "http://jira.com", prReviewed: false, state: sprintbot.IssueStateReadyForQA},
			},
			Err:            nil,
			IssueRepoCalls: map[string]int{"SaveNext": 1},
		},
		{
			Name:        "should not find and return pr issues for ignored repo",
			ExpectError: false,
			NumIssues:   0,
			Issues: []mockIssue{
				mockIssue{pr: "https://gitlab.cee.redhat.com/red-hat-mobile-application-platform-documentation/RHMAPDocsNG/merge_requests/181", link: "http://jira.com", prReviewed: false, state: sprintbot.IssueStateClosed},
			},
			Err:            nil,
			IssueRepoCalls: map[string]int{"SaveNext": 1},
		},
		{
			Name:        "should find and return ready for qe  issues",
			ExpectError: false,
			NumIssues:   1,
			Issues: []mockIssue{
				mockIssue{pr: "github.com/some/pr/20", link: "http://jira.com", prReviewed: true, state: sprintbot.IssueStateReadyForQA},
			},
			Err:            nil,
			IssueRepoCalls: map[string]int{"SaveNext": 1},
		},
		{
			Name:        "should find and return Open  issues",
			ExpectError: false,
			NumIssues:   1,
			Issues: []mockIssue{
				mockIssue{pr: "github.com/some/pr/20", link: "http://jira.com", prReviewed: true, state: sprintbot.IssueStateOpen},
			},
			Err:            nil,
			IssueRepoCalls: map[string]int{"SaveNext": 1},
		},
		{
			Name:        "should add comment to issues where prs are reviewed but still in PR Sent state",
			ExpectError: false,
			NumIssues:   1,
			Issues: []mockIssue{
				mockIssue{pr: "github.com/some/pr/20", link: "http://jira.com", prReviewed: true, state: sprintbot.IssueStateOpen},
				mockIssue{pr: "github.com/some/pr/20", link: "http://jira.com", prReviewed: true, state: sprintbot.IssueStatePRSent},
			},
			Err:                    nil,
			IssueRepoCalls:         map[string]int{"SaveNext": 1, "FindCommentOnIssue": 1, "SaveCommented": 1},
			IssueEditorFinderCalls: map[string]int{"AddComment": 1},
		},
	}

	for _, tc := range cases {
		isf := newMockIssueEditorFinder(tc.Issues, tc.Err)
		sprintFinder := newMockSprintFinder(nil, nil)
		rc := &mockRepoChecker{issues: tc.Issues}
		s := sprintbot.Sprint{Name: "testSprint", Board: "testBoard"}
		logger := logrus.StandardLogger()
		issueRepo := newMockIssueRepo(nil, "", nil, nil)
		service := sprint.NewService(isf, sprintFinder, rc, issueRepo, &s, logger)
		service.IgnoredRepos = []string{"RHMAPDocsNG", "fhcap", "fh-openshift-templates", "fh-core-openshift-templates"}
		t.Run(tc.Name, func(t *testing.T) {
			err := service.Sync()
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none ")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect an error but got one  %s ", err)
			}
			if tc.IssueRepoCalls != nil {
				for k, v := range tc.IssueRepoCalls {
					if issueRepo.Called(k) != v {
						t.Fatalf("expected %s to be called %v time but got %v", k, v, issueRepo.Called(k))
					}
				}
			}
			if tc.IssueEditorFinderCalls != nil {
				for k, v := range tc.IssueEditorFinderCalls {
					if isf.Called(k) != v {
						t.Fatalf("expected %s to be called %v time but got %v", k, v, isf.Called(k))
					}
				}
			}
		})
	}
}

//TODO needs a more advanced test
func TestSprintStatus(t *testing.T) {
	cases := []struct {
		Name         string
		Err          error
		ExpectError  bool
		OpenPoints   int
		ClosedPoints int
		Issues       map[string][]sprintbot.IssueState
		Sprint       *sprintbot.JiraSprint
		TimeCalc     func(start, end string) (*sprintbot.Time, error)
	}{
		{
			Name:         "test status happy",
			Err:          nil,
			ExpectError:  false,
			OpenPoints:   10,
			ClosedPoints: 5,
			Issues: map[string][]sprintbot.IssueState{
				sprintbot.IssueStateOpen: []sprintbot.IssueState{
					mockIssue{state: sprintbot.IssueStateOpen, storyPoints: 5},
					mockIssue{state: sprintbot.IssueStateOpen, storyPoints: 5},
				},
				sprintbot.IssueStateClosed: []sprintbot.IssueState{
					mockIssue{state: sprintbot.IssueStateClosed, storyPoints: 5},
				},
			},
			Sprint: &sprintbot.JiraSprint{},
			TimeCalc: func(start, end string) (*sprintbot.Time, error) {
				return &sprintbot.Time{
					DaysGone: 2,
					DaysLeft: 2,
				}, nil
			},
		},
	}

	for _, tc := range cases {
		isf := newMockIssueEditorFinder(nil, tc.Err)
		sprintFinder := newMockSprintFinder(tc.Sprint, tc.Err)
		rc := &mockRepoChecker{issues: nil}
		s := sprintbot.Sprint{Name: "testSprint", Board: "testBoard"}
		logger := logrus.StandardLogger()
		issueRepo := newMockIssueRepo(nil, "", nil, tc.Issues)
		service := sprint.NewService(isf, sprintFinder, rc, issueRepo, &s, logger)
		service.TimeCalc = tc.TimeCalc
		service.IgnoredRepos = []string{"RHMAPDocsNG", "fhcap", "fh-openshift-templates", "fh-core-openshift-templates"}
		t.Run(tc.Name, func(t *testing.T) {
			status, err := service.Status()
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect error but got one %s ", err)
			}
			if !tc.ExpectError && status == nil {
				t.Fatalf("status should not have been nil")
			}
			if tc.ClosedPoints != status.PointsCompleted {
				t.Fatalf("expected the points closed to match got %v and %v ", tc.ClosedPoints, status.PointsCompleted)
			}
			if tc.OpenPoints != status.PointsRemaining {
				t.Fatalf("expected the points open to match got %v and %v ", tc.OpenPoints, status.PointsRemaining)
			}
		})
	}
}
