package usecase

import (
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
func (ch *Chat) Handle(cmd sprintbot.ChatCMD) (*sprintbot.ChatResponse, error) {
	switch cmd.Action() {
	case sprint.CommandNext:
		ch.sprintService.Next()
	default:
		return nil, &sprintbot.ErrUnkownCMD{Message: "the command " + cmd.Action() + " is not something I can do"}
	}
	return nil, nil
}
