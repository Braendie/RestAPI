package handlers

import "github.com/julienschmidt/httprouter"

// Handler is a interface for our struct handler
type Handler interface {
	Register(router *httprouter.Router) 
}