package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

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
	switch req.Method {

	// Handle get requests
	case http.MethodGet:
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

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
