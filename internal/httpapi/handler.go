package httpapi

import (
	"encoding/json"
	domain "micro-vuln-scanner/internal/domain"
	store "micro-vuln-scanner/internal/store"
	"net/http"
)

type HttpHandler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *HttpHandler {
	return &HttpHandler{
		store: store,
	}
}

func (handler *HttpHandler) GetVulnerabilities(w http.ResponseWriter, request *http.Request) {
	var vulnerabilities []domain.Vulnerability
	w.Header().Set("Content-Type", "application/json")
	if request.URL.Query().Get("severity") != "" {
		severity, err := domain.ParseSeverity(request.URL.Query().Get("severity"))
		if err != nil {
			http.Error(w, "{\"error\": \"Invalid severity value\"}", http.StatusBadRequest)
			return
		}
		 vulnerabilities = handler.store.GetBySeverity(severity)
	} else {
	vulnerabilities = handler.store.GetAll()
	}

	err := json.NewEncoder(w).Encode(vulnerabilities)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}