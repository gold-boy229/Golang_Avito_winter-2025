package handlers

import (
	"MerchShop/model"
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, model.ErrorResponse{Errors: message})
}

func respondBadRequest(w http.ResponseWriter, message string) {
	respondError(w, http.StatusBadRequest, message)
}

func respondUnauthorized(w http.ResponseWriter, message string) {
	respondError(w, http.StatusUnauthorized, message)
}

func respondNotFound(w http.ResponseWriter, message string) {
	respondError(w, http.StatusNotFound, message)
}

func respondMethodNotAllowed(w http.ResponseWriter, message string) {
	respondError(w, http.StatusNotFound, message)
}

func respondInternalServerError(w http.ResponseWriter, message string) {
	respondError(w, http.StatusInternalServerError, message)
}
