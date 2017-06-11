package sprintbot

type Chatter interface {
	Chat() error
}

type IssueFinder interface {
	FindUnresolvedOnBoard(boardName, sprint string) (*IssueList, error)
	IssueHost() string
}

type IssueRepo interface {
	SaveNext(next *NextIssues) error
	FindNext() (*NextIssues, error)
}

type IssueState interface {
	PRS() []string
	RemovePR(pr string)
	Link(host string) string
	State() string
	Description() string
}

type RepoChecker interface {
	PRReviewed(prURL string) (bool, error)
	Repo(pURL string) string
}

type ChatCMD interface {
	Action() string
	User() string
	AuthToken() string
}
