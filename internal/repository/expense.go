package repository

import (
	"database/sql"
	"errors"
	"expense-tracker/internal/model"
)

var ErrNotFound = errors.New("expense not found")

type ExpenseRepo struct {
	DB *sql.DB
}

func NewExpenseRepo(db *sql.DB) *ExpenseRepo {
	return &ExpenseRepo{DB: db}
}

func (r *ExpenseRepo) Patch(id int64, e *model.Expense) (*model.Expense, error) {
	var updated model.Expense

	err := r.DB.QueryRow(`
	UPDATE expenses SET
	amount = COALESCE($1, amount),
	note = COALESCE($2, not)
	WHERE id=$3
	RETURNING id, date, amount, note`,
		updated.Amount, updated.Note, id).Scan(
		&updated.ID, &updated.Date, &updated.Amount, &updated.Note)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *ExpenseRepo) GetById(id int64, e *model.Expense) (*model.Expense, error) {
	var eId model.Expense
	err := r.DB.QueryRow(`
	SELECT id, date, amount, note
	FROM expenses
	WHERE id=$1`, id).Scan(
		&eId.ID, &eId.Date, &eId.Amount, &eId.Note)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &eId, nil

}

func (r *ExpenseRepo) Delete(id int64) error {
	result, err := r.DB.Exec(`DELETE FROM expenses WHERE id=$1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	return nil
}
