package web

import "net/http"

type ChatHandler struct{}

// IncomingMessage handles messages coming from the chat server directed at the sprint bot
func (ch *ChatHandler) IncomingMessage(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("test"))
}
