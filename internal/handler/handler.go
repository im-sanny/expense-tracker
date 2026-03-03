package handler

import (
	"encoding/json"
	"errors"
	"expense-tracker/internal/model"
	"expense-tracker/internal/repository"
	"expense-tracker/internal/service"
	"expense-tracker/pkg/apperrors"
	"expense-tracker/pkg/response"
	"log"
	"net/http"
	"strconv"
)

type ExpenseHandler struct {
	Service service.ExpenseServiceInterface
}

func NewHandler(service service.ExpenseServiceInterface) *ExpenseHandler {
	return &ExpenseHandler{Service: service}
}

// =======GET=======
func (h *ExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {

	params, err := ParseExpenseQuery(r)
	if err != nil {
		response.WriteBadRequest(w, err)
		return
	}
	log.Printf("Parsed params: %+v", params)

	filter := repository.ExpenseFilter{
		Min:    params.Min,
		Max:    params.Max,
		From:   params.From,
		To:     params.To,
		Search: params.Search,
	}

	result, err := h.Service.Get(r.Context(), params.Page, params.Limit, filter)
	if err != nil {
		log.Printf("server error: %v", err)
		response.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	if err := response.WriteSuccess(w, result); err != nil {
		log.Printf("failed to write response: %v", err)
		return
	}
}

// =======GETbyID=======
func (h *ExpenseHandler) GetById(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	expense, err := h.Service.GetById(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get expense", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(expense); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// =======POST=======
func (h *ExpenseHandler) Post(w http.ResponseWriter, r *http.Request) {
	var e model.Expense

	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	create, err := h.Service.Post(r.Context(), &e)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "expense not created", http.StatusNoContent)
			return
		}
		log.Println(err)
		http.Error(w, "failed to insert data in expenses", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(create)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// =======PUT=======
func (h *ExpenseHandler) Put(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var e model.Expense

	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	updated, err := h.Service.Put(r.Context(), int64(id), &e)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update expense", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(updated); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

// =======PATCH=======
func (h *ExpenseHandler) Patch(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var e model.Expense
	if err = json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	updated, err := h.Service.Patch(r.Context(), int64(id), &e)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update expense", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(updated); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// "I use pointers so that only the fields actually sent in the PATCH request get updated. Fields not included are left unchanged."

// =======DELETE=======
func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	err = h.Service.Delete(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete expense", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
