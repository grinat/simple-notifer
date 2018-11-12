package main

import (
	"net/http"
)

type Router struct {
	cfg Params
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/send-message":
		HandleSendMessage(w, r, router.cfg)
	default:
		HandleHelp(w, router.cfg)
	}
}