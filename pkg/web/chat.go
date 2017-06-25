package web

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/maleck13/sprintbot/pkg/sprintbot/usecase"
)

type ChatHandler struct {
	chat      *usecase.Chat
	logger    *logrus.Logger
	authToken string
}

func NewChatHandler(chatUseCase *usecase.Chat, logger *logrus.Logger, auth string) *ChatHandler {
	return &ChatHandler{
		chat:      chatUseCase,
		logger:    logger,
		authToken: auth,
	}
}

// IncomingMessage handles messages coming from the chat server directed at the sprint bot
func (ch *ChatHandler) IncomingMessage(rw http.ResponseWriter, req *http.Request) {
	var (
		decoder = json.NewDecoder(req.Body)
		encoder = json.NewEncoder(rw)
	)
	source := req.URL.Query().Get("source")
	if "" == source {
		http.Error(rw, "uknown source in query string. It cannot be empty", http.StatusBadRequest)
		return
	}
	var chatCMD sprintbot.ChatCMD
	switch source {
	case "rocket":
		chatCMD = &sprintbot.RocketChatCmd{}
	}
	if err := decoder.Decode(chatCMD); err != nil {
		http.Error(rw, "failed decoding incoming message. "+err.Error(), http.StatusBadRequest)
		return
	}
	if ch.authToken != chatCMD.AuthToken() {
		http.Error(rw, "not authorised", http.StatusUnauthorized)
		return
	}
	chatResp, err := ch.chat.Handle(chatCMD)
	if err != nil {
		ch.handleChatError(err, rw)
		return
	}
	if err := encoder.Encode(chatResp); err != nil {
		http.Error(rw, "failed to encode response to chat cmd "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ch *ChatHandler) handleChatError(err error, rw http.ResponseWriter) {
	ch.logger.Error(err)
	switch err.(type) {
	case *sprintbot.ErrUnkownCMD, *sprintbot.ErrInvalid:
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(rw, "unexpected error "+err.Error(), http.StatusInternalServerError)
		return
	}
}
