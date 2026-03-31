package httpapi

import (
	"encoding/json"
	service "micro-vuln-scanner/internal/service"
	store "micro-vuln-scanner/internal/store"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message})
}

type HttpHandler struct {
	Store *store.Store
}

func NewHandler(store *store.Store) *HttpHandler {
	return &HttpHandler{
		Store: store,
	}
}

func (handler *HttpHandler) GetVulnerabilities(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	severity := request.URL.Query().Get("severity")
	vulnerabilities, err := service.GetVulnerabilities(handler.Store, severity)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid severity parameter")
		return
	}

	if err := json.NewEncoder(w).Encode(vulnerabilities); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to encode vulnerabilities")
		return
	}

}