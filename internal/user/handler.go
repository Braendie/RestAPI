package user

import (
	"net/http"

	"github.com/Braendie/RestAPI/internal/handlers"
	"github.com/julienschmidt/httprouter"
)

// endpoints
const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)


type handler struct {
}

// NewHandler creates new handler and return its
func NewHandler() handlers.Handler {
	return &handler{}
}

// Register registers all the endpoints of this server
func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.GetList)
	router.GET(userURL, h.GetUserByUUID)
	router.POST(usersURL, h.CreateUser)
	router.PUT(userURL, h.UpdateUser)
	router.PATCH(userURL, h.PartiallyUpdateUser)
	router.DELETE(userURL, h.DeleteUser)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is list of users"))
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is user by uuid"))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is create user"))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is update user"))
}

func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is partially update user"))
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is delete user"))
}
