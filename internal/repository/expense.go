package repository

import (
	"database/sql"
	"errors"
	"expense-tracker/internal/model"
	"time"
)

var ErrNotFound = errors.New("expense not found")

type ExpenseRepo struct {
	DB *sql.DB
}

func NewExpenseRepo(db *sql.DB) *ExpenseRepo {
	return &ExpenseRepo{DB: db}
}

func (r *ExpenseRepo) Post(expense *model.Expense) (*model.Expense, error) {
	var created model.Expense
	expense.Date = time.Now()

	err := r.DB.QueryRow(`
	INSERT INTO
	expenses(date, amount, note)
	VALUES($1, $2, $3)
	RETURNING id, date, amount, note`,
		expense.Date, expense.Amount, expense.Note).Scan(
		&created.ID, &created.Date, &created.Amount, &created.Note)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *ExpenseRepo) Get() ([]model.Expense, error) {
	var expenses []model.Expense
	rows, err := r.DB.Query(`SELECT id, date, amount, note FROM expenses`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var expense model.Expense
		if err := rows.Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.Note); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *ExpenseRepo) GetById(id int64) (*model.Expense, error) {
	var expense model.Expense
	err := r.DB.QueryRow(`
	SELECT id, date, amount, note
	FROM expenses
	WHERE id=$1`, id).Scan(
		&expense.ID, &expense.Date, &expense.Amount, &expense.Note)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *ExpenseRepo) Put(id int64, expense *model.Expense) (*model.Expense, error) {
	expense.Date = time.Now()
	err := r.DB.QueryRow(`
	UPDATE expenses SET
	amount=$1, note=$2 WHERE id=$3
	RETURNING id, date, amount, note`,
		expense.Amount, expense.Note, id).Scan(
		&expense.ID, &expense.Date, &expense.Amount, &expense.Note)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return expense, nil
}

func (r *ExpenseRepo) Patch(id int64, update *model.Expense) (*model.Expense, error) {
	var expense model.Expense

	err := r.DB.QueryRow(`
	UPDATE expenses SET
	amount = COALESCE($1, amount),
	note = COALESCE($2, note)
	WHERE id=$3
	RETURNING id, date, amount, note`,
		update.Amount, update.Note, id).Scan(
		&expense.ID, &expense.Date, &expense.Amount, &expense.Note)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &expense, nil
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
