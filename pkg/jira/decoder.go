package jira

import (
	"encoding/json"

	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/pkg/errors"
)

type Decoder struct {
}

func (d *Decoder) Decode(data []byte) ([]sprintbot.IssueState, error) {

	var issues = []*sprintbot.JiraIssue{}
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal into Jira issues")
	}
	var issueState = make([]sprintbot.IssueState, len(issues))
	for i, is := range issues {
		issueState[i] = is
	}
	return issueState, nil
}
