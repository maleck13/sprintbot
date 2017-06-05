package sprintbot

// Chat handles the chat usecase
type Chat struct {
}

// Handle will take a chat command and process it returning a chat response
func (ch *Chat) Handle(cmd *ChatCmd) (*ChatResponse, error) {
	return nil, nil
}
