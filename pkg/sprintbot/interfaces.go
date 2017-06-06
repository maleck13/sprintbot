package sprintbot

type Chatter interface {
	Chat() error
}

type IssueFinder interface {
	FindUnresolved() (*IssueList, error)
}

type IssueState interface {
	PR() string
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
