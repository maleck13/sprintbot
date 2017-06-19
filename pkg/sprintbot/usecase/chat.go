package usecase

import (
	"fmt"

	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/maleck13/sprintbot/pkg/sprintbot/sprint"
	"github.com/pkg/errors"
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
func (ch *Chat) Handle(cmd sprintbot.ChatCMD) (*sprintbot.ChatResponse, error) {
	fmt.Println("handling cmd action ", cmd.Action())
	var err error
	var resp *sprintbot.ChatResponse
	switch cmd.Action() {
	case sprint.CommandNext:
		next, err := ch.sprintService.Next()
		if err != nil {
			return nil, errors.Wrap(err, "handle chat command failed after calling next")
		}
		resp = &sprintbot.ChatResponse{Message: fmt.Sprintf("@%s %s", cmd.User(), next.Message), Data: next.Issues, CMD: sprint.CommandNext}
	case sprint.CommandStatus:
		status, err := ch.sprintService.Status()
		if err != nil {
			return nil, errors.Wrap(err, "handle chat command failed after calling status")
		}
		resp = &sprintbot.ChatResponse{Message: fmt.Sprintf("@%s %s", cmd.User(), "sprint status"), Data: status, CMD: sprint.CommandStatus}
	default:
		return nil, &sprintbot.ErrUnkownCMD{Message: "the command " + cmd.Action() + " is not something I can do"}
	}
	return resp, err
}
