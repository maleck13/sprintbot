package usecase

import (
	"fmt"

	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/maleck13/sprintbot/pkg/sprintbot/sprint"
)

// Chat handles the chat usecase
type Chat struct {
	sprintService *sprint.Service
}

// NewChat returns a chat usecase handler
func NewChat(sprintService *sprint.Service) *Chat {
	return &Chat{
		sprintService: sprintService,
	}
}

// Handle will take a chat command and process it returning a chat response
func (ch *Chat) Handle(cmd sprintbot.ChatCMD) (*sprintbot.NextIssues, error) {
	fmt.Println("handling cmd action ", cmd.Action())
	var err error
	var resp *sprintbot.NextIssues
	switch cmd.Action() {
	case sprint.CommandNext:
		resp, err = ch.sprintService.Next()
	default:
		return nil, &sprintbot.ErrUnkownCMD{Message: "the command " + cmd.Action() + " is not something I can do"}
	}
	resp.Message = fmt.Sprintf("@%s %s", cmd.User(), resp.Message)
	return resp, err
}
