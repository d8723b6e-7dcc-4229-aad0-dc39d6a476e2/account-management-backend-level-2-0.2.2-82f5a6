package account

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type accountHandler struct {
	BalanceRouter *mux.Router
	AmountRouter  *mux.Router
	service       Service
}

// NewHandler returns new account handler
func NewHandler(s Service) *accountHandler {
	handler := &accountHandler{
		service: s,
	}
	handler.BalanceRouter = balanceRoutes(handler)
	handler.AmountRouter = amountRoutes(handler)

	return handler
}

func balanceRoutes(handler *accountHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/balance/{account_id}", handler.accountBalanceHandlerFunc).Methods(http.MethodGet)
	return r
}

func amountRoutes(handler *accountHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/amount", handler.amountHandlerFunc).Methods(http.MethodPost)
	return r
}

func (h *accountHandler) accountBalanceHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["account_id"]

	balance, err := h.service.GetAccountBalance(r.Context(), accountID)
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := balanceResponse{balance}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

type balanceResponse struct {
	Balance int `json:"balance"`
}

type amountRequest struct {
	AccountID *string `json:"account_id"`
	Amount    *int    `json:"amount"`
}

func (h *accountHandler) amountHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	var req amountRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Ensure that account id and amount are provided in request
	if req.AccountID == nil || req.Amount == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Ensure that transaction id is provided in header
	transactionID := r.Header.Get("Transaction-ID")
	if transactionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.service.UpdateAccountBalance(r.Context(), *req.AccountID, transactionID, *req.Amount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func decodeJSON(r io.Reader, input interface{}) error {
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()

	e := d.Decode(input)
	if e != nil {
		return e
	}

	return nil
}
