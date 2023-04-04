package transaction

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type accountHandler struct {
	Router  *mux.Router
	service Service
}

func NewHandler(s Service) *accountHandler {
	handler := &accountHandler{
		service: s,
	}
	handler.Router = route(handler)
	return handler
}

func route(handler *accountHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/transaction/{transaction_id}", handler.transactionDetailsHandlerFunc).Methods(http.MethodGet)
	return r
}

func (h *accountHandler) transactionDetailsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transaction_id"]

	transaction, err := h.service.GetTransactionByID(r.Context(), transactionID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if transaction == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}
