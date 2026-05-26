package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"k-microserv-kuber.com/users"
	"k-microserv-kuber.com/users/controller"
	"k-microserv-kuber.com/users/database"
)

type Handler struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// Handle requests on /user endpoint
func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {

	id := req.FormValue("userId")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {

	// Handle get requests
	case http.MethodGet:

		ctx := req.Context()
		u, err := h.ctrl.GetUser(ctx, userId)
		if err != nil && errors.Is(err, database.ErrNotFound) {
			log.Println("User not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("User get error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(u); err != nil {
			log.Printf("User get response encode error: %v\n", err)
		}

	case http.MethodPost:
		email := req.FormValue("email")
		if email == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := req.Context()
		newUser := &users.User{UserId: userId, Email: email}
		if err := h.ctrl.PutUser(ctx, userId, newUser); err != nil {
			if errors.Is(err, database.UserAlreadyExists) {
				log.Println("User already exists")
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				log.Printf("User creation error: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
