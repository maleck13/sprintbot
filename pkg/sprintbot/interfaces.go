package sprintbot

type Chatter interface {
	Chat() error
}

type IssueFinder interface {
	IssuesForBoard(boardName, sprint string) (*IssueList, error)
	IssueHost() string
}

type SprintFinder interface {
	SprintForBoard(sprintName, boardName string) (*JiraSprint, error)
}

type IssueEditorFinder interface {
	IssueFinder
	AddComment(id, c string) error
}

type IssueDecoder interface {
	Decode(data []byte) ([]IssueState, error)
}

type IssueRepo interface {
	SaveNext(next *NextIssues) error
	FindNext() (*NextIssues, error)
	SaveCommented(id string, commentID string) error
	FindCommentOnIssue(id string, commentID string) (string, error)
	SaveIssuesInState(state string, issues []IssueState) error
	FindIssuesInState(state string) ([]IssueState, error)
	FindIssuesNotInStates(states []string) ([]IssueState, error)
}

type IssueState interface {
	ID() string
	PRS() []string
	RemovePR(pr string)
	Link(host string) string
	State() string
	Description() string
	StoryPoints() int
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
