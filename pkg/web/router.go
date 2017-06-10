package web

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/maleck13/sprintbot/pkg/sprintbot/usecase"
	"github.com/rs/cors"
)

// BuildRouter is the main place we build the mux router
func BuildRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return r
}

// BuildHTTPHandler constructs a http.Handler it is also where common middleware is added via negroni
func BuildHTTPHandler(r *mux.Router) http.Handler {
	//recovery middleware for any panics in the handlers
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	//add middleware for all routes
	n := negroni.New(recovery)
	n.Use(cors.New(
		cors.Options{
			AllowedOrigins: []string{"*"},
		},
	))
	// set up sys routes
	n.UseHandler(r)
	return n
}

func ChatRoute(r *mux.Router, cu *usecase.Chat, logger *logrus.Logger, auth string) {
	chatHandler := NewChatHandler(cu, logger, auth)
	r.HandleFunc("/chat/message", chatHandler.IncomingMessage).Methods("POST")
}
