package sprint

import (
	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/pkg/errors"
)

const (
	// CommandNext is the command sent from the chat server to find out whats next in the sprint
	CommandNext = "next"
)

// Service handles the buisness logic around sprint actions
type Service struct {
	issueFinder  sprintbot.IssueFinder
	repoChecker  sprintbot.RepoChecker
	boardName    string
	sprintName   string
	IgnoredRepos []string
}

// NewService returns a new  instance of the sprint service
func NewService(issueFinder sprintbot.IssueFinder, repoChecker sprintbot.RepoChecker, sp *sprintbot.Sprint) *Service {
	return &Service{
		issueFinder:  issueFinder,
		repoChecker:  repoChecker,
		sprintName:   sp.Name,
		boardName:    sp.Board,
		IgnoredRepos: []string{}, // this doesn't feel quite right
	}
}

func (s *Service) awaitingPrReview(issues []sprintbot.IssueState) ([]sprintbot.IssueState, error) {
	var awaitingPR = []sprintbot.IssueState{}
	for _, i := range issues {
		if len(i.PRS()) == 0 {
			continue
		}
		for _, pr := range i.PRS() {
			shouldIgnore := false
			repo := s.repoChecker.Repo(pr)
			for _, ignore := range s.IgnoredRepos {
				if ignore == repo {
					shouldIgnore = true
				}
			}
			if shouldIgnore {
				continue
			}
			is, err := s.repoChecker.PRReviewed(pr)
			if err != nil {
				return nil, errors.Wrap(err, " sprint service next failed checking if PR is reviewed ")
			}
			if !is {
				awaitingPR = append(awaitingPR, i)
			}
		}
	}
	return awaitingPR, nil
}

func (s *Service) buildIssues(is []sprintbot.IssueState) []*sprintbot.Issue {
	next := make([]*sprintbot.Issue, len(is))
	for i, ni := range is {
		next[i] = &sprintbot.Issue{
			Link:        ni.Link(s.issueFinder.IssueHost()),
			PRs:         ni.PRS(),
			Description: ni.Description(),
		}
	}
	return next
}

// Next will look at the current sprint and figure out what should be the next action
func (s *Service) Next() (*sprintbot.NextIssues, error) {
	if s.boardName == "" || s.sprintName == "" {
		return nil, &sprintbot.ErrInvalid{Message: "expected a board name and sprint name but did not get them boardName: " + s.boardName + " sprintName: " + s.boardName}
	}
	var next = []*sprintbot.Issue{}
	// retreive  all non closed issues
	issueList, err := s.issueFinder.FindUnresolvedOnBoard(s.boardName, s.sprintName)
	if err != nil {
		return nil, errors.Wrap(err, "sprint service finding next failed to find unresolved issues ")
	}
	issues := issueList.Issues()
	if len(issues) == 0 {
		return &sprintbot.NextIssues{Issues: s.buildIssues(issues), Message: "nothing to do?"}, nil
	}
	// check for PR sent and not reviewed PRs
	awaitingPR, err := s.awaitingPrReview(issues)
	if err != nil {
		return nil, err
	}
	if len(awaitingPR) > 0 {
		return &sprintbot.NextIssues{Issues: s.buildIssues(awaitingPR), Message: "There are some outstanding PRs on these issues."}, nil
	}
	// if none check for ready for qe but not qe in progress or PR sent where the pr has been reviewed
	rqaissueList := issueList.FindInState(sprintbot.IssueStateReadyForQA)
	issues = rqaissueList.Issues()
	if len(issues) > 0 {
		return &sprintbot.NextIssues{Issues: s.buildIssues(issues), Message: "The following issues are in a " + sprintbot.IssueStateReadyForQA + " state."}, nil
	}
	//if none look for the next available issue
	issueList = issueList.FindInState(sprintbot.IssueStateOpen)
	issues = issueList.Issues()
	if len(issues) > 0 {
		return &sprintbot.NextIssues{Issues: s.buildIssues(issues), Message: "The following issues are in a " + sprintbot.IssueStateOpen + " state"}, nil
	}
	return &sprintbot.NextIssues{Issues: next, Message: "Looks like there is nothing to do?"}, nil
}
