package sprintbot

type Chatter interface {
	Chat() error
}

type IssueFinder interface {
	FindUnresolvedOnBoard(boardName, sprint string) (*IssueList, error)
}

type IssueState interface {
	PRS() []string
	Link() string
	State() string
}

type RepoChecker interface {
	PRReviewed(prURL string) (bool, error)
}

type ChatCMD interface {
	Action() string
	User() string
}
