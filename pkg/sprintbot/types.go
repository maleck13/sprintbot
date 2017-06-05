package sprintbot

type ChatCmd struct{}
type ChatResponse struct{}

type NextIssues struct {
	Message string
	Issues  []IssueState
}

type IssueList struct {
	issues []IssueState
}

func NewIssueList(issues []IssueState) *IssueList {
	return &IssueList{
		issues: issues,
	}
}

func (il *IssueList) Issues() []IssueState {
	return il.issues
}

func (il *IssueList) FindInState(state string) *IssueList {
	var issues = []IssueState{}
	for _, i := range il.issues {
		if i.State() == state {
			issues = append(issues, i)
		}
	}
	return &IssueList{issues: issues}
}

const (
	IssueStateReadyForQA = "Ready for QA"
	IssueStateOpen       = "Open"
)
