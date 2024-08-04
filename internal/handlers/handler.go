package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kylejryan/mocument/internal/db"
)

type Handler struct {
	docDB db.DocumentDB
}

func NewHandler(docDB db.DocumentDB) *Handler {
	return &Handler{
		docDB: docDB,
	}
}

func (h *Handler) HandleRequests() {
	http.HandleFunc("/insert", h.insertDocument)
	http.HandleFunc("/find", h.findDocument)
	http.HandleFunc("/update", h.updateDocuments)
	// http.HandleFunc("/delete", h.deleteDocument)
	http.ListenAndServe(":8080", nil)
}

func (h *Handler) insertDocument(w http.ResponseWriter, r *http.Request) {
	var doc map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.docDB.InsertDocument("collection", doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) findDocument(w http.ResponseWriter, r *http.Request) {
	filterParam := r.URL.Query().Get("filter")
	var filter map[string]interface{}
	if filterParam != "" {
		if err := json.Unmarshal([]byte(filterParam), &filter); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	result, err := h.docDB.FindDocument("collection", filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) updateDocuments(w http.ResponseWriter, r *http.Request) {
	var updateData struct {
		Filter map[string]interface{} `json:"filter"`
		Update map[string]interface{} `json:"update"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.docDB.UpdateMany("collection", updateData.Filter, updateData.Update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
