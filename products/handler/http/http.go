package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"k-microserv-kuber.com/products"
	"k-microserv-kuber.com/products/controller"
	"k-microserv-kuber.com/products/database"
)

type Handler struct {
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("productId")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		ctx := req.Context()
		p, err := h.ctrl.GetProduct(ctx, productId)
		if err != nil && errors.Is(err, database.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Product get error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(p); err != nil {
			log.Printf("Product get response encode error: %v\n", err)
		}
	case http.MethodPost:
		productName := req.FormValue("name")
		if productName == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		priceString := req.FormValue("price")
		if priceString == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(priceString, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := req.Context()
		newProduct := &products.Product{
			ProductID:   productId,
			ProductName: productName,
			Price:       price,
		}
		if err := h.ctrl.InsertProduct(ctx, productId, newProduct); err != nil {
			if errors.Is(err, database.ProductAlreadyExists) {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				log.Printf("Product creation error: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
