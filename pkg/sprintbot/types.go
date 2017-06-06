package sprintbot

import (
	"fmt"
	"strings"
)

type RocketChatCmd struct {
	Bot         bool   `json:"bot"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	IsEdited    bool   `json:"isEdited"`
	MessageID   string `json:"message_id"`
	Text        string `json:"text"`
	Timestamp   string `json:"timestamp"`
	Token       string `json:"token"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
}

func (rcmd *RocketChatCmd) Action() string {
	return strings.Replace(rcmd.Text, "sprintbot ", "", 1)
}

func (rcmd *RocketChatCmd) User() string {
	return rcmd.UserName
}

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

type ErrUnkownCMD struct {
	Message string
}

func (e *ErrUnkownCMD) Error() string {
	return fmt.Sprintf("Error unknown cmd: %s", e.Message)
}
