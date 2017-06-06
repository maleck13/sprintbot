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
	issueFinder sprintbot.IssueFinder
	repoChecker sprintbot.RepoChecker
}

// NewService returns a new  instance of the sprint service
func NewService(issueFinder sprintbot.IssueFinder, repoChecker sprintbot.RepoChecker) *Service {
	return &Service{
		issueFinder: issueFinder,
		repoChecker: repoChecker,
	}
}

// Next will look at the current sprint and figure out what should be the next action
func (s *Service) Next() (*sprintbot.NextIssues, error) {
	// retreive  all non closed issues
	issueList, err := s.issueFinder.FindUnresolved()
	if err != nil {
		return nil, errors.Wrap(err, "sprint service finding next failed to find unresolved issues ")
	}
	issues := issueList.Issues()
	if len(issues) == 0 {
		return &sprintbot.NextIssues{Issues: issues, Message: "nothing to do?"}, nil
	}
	// check for PR sent and not reviewed PRs
	var awaitingPR = []sprintbot.IssueState{}
	for _, i := range issues {
		if "" != i.PR() {
			is, err := s.repoChecker.PRReviewed(i.PR())
			if err != nil {
				return nil, errors.Wrap(err, " sprint service next failed checking if PR is reviewed ")
			}
			if !is {
				awaitingPR = append(awaitingPR, i)
			}
		}
	}
	if len(awaitingPR) > 0 {
		return &sprintbot.NextIssues{Issues: awaitingPR, Message: "There are some outstanding PRs on these issues."}, nil
	}
	// if none check for ready for qe but not qe in progress or PR sent where the pr has been reviewed
	rqaissueList := issueList.FindInState(sprintbot.IssueStateReadyForQA)
	issues = rqaissueList.Issues()
	if len(issues) > 0 {
		return &sprintbot.NextIssues{Issues: issues, Message: "The following issues are in " + sprintbot.IssueStateReadyForQA}, nil
	}
	//if none look for the next available issue
	issueList = issueList.FindInState(sprintbot.IssueStateOpen)
	issues = issueList.Issues()
	if len(issues) > 0 {
		return &sprintbot.NextIssues{Issues: issues, Message: "The following issues are in " + sprintbot.IssueStateOpen}, nil
	}
	return &sprintbot.NextIssues{Issues: issues, Message: "Looks like there is nothing to do?"}, nil
}
