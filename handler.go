package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// MyHandler структура для обработки HTTP-запросов.
type MyHandler struct {
	service Service
}

// NewMyHandler создает новый экземпляр MyHandler.
func NewMyHandler(service Service) *MyHandler {
	return &MyHandler{service: service}
}

// SearchHandler обрабатывает поисковые запросы.
func (h *MyHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	vacancies, err := h.service.SearchVacancy(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching for vacancies: %v", err), http.StatusInternalServerError)
		return
	}

	// Добавление обработки ошибок при записи в лог или базу данных.
	if err := h.service.SaveSearchHistory(query); err != nil {
		fmt.Printf("Error saving search history: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vacancies)
}

// GetHandler обрабатывает запросы на получение вакансии.
func (h *MyHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	vacancyID, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "Vacancy ID is required", http.StatusBadRequest)
		return
	}

	vacancy, err := h.service.GetVacancy(vacancyID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting vacancy: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, vacancy)
}

// ListHandler обрабатывает запросы на получение списка вакансий.
func (h *MyHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	vacancies, err := h.service.ListVacancies()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting list of vacancies: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, vacancies)
}

// DeleteHandler обрабатывает запросы на удаление вакансии.
func (h *MyHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vacancyID, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "Vacancy ID is required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteVacancy(vacancyID); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting vacancy: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListSearchHistoryHandler обрабатывает запросы на получение истории поиска.
func (h *MyHandler) ListSearchHistoryHandler(w http.ResponseWriter, r *http.Request) {
	history, err := h.service.ListSearchHistory()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting search history: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, history)
}

// jsonResponse упрощает отправку JSON-ответов.
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func (h *MyHandler) DeleteSearchHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	historyID, ok := vars["id"]
	if !ok {
		http.Error(w, "History ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(historyID)
	if err != nil {
		http.Error(w, "Invalid History ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteSearchHistory(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting search history: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
