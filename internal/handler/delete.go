package handler

import (
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Exec(`DELETE FROM expenses WHERE id=$1`, id)
	if err != nil {
		log.Println(err)
		return
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "data not found", http.StatusNotFound)
		return
	}

}
