package repository

import (
	"database/sql"
	"expense-tracker/internal/model"
	"expense-tracker/pkg/errors"
	"time"
)

type ExpenseFilter struct {
	Min  int       `json:"min"`
	Max  int       `json:"max"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type ExpenseRepoInterface interface {
	Get(offset, limit int, f ExpenseFilter) ([]model.Expense, error)
	GetById(id int64) (*model.Expense, error)
	Post(e *model.Expense) (*model.Expense, error)
	Put(id int64, e *model.Expense) (*model.Expense, error)
	Patch(id int64, e *model.Expense) (*model.Expense, error)
	Delete(id int64) error
	Count(f ExpenseFilter) (int, error)
}

type ExpenseRepo struct {
	DB *sql.DB
}

func NewExpenseRepo(db *sql.DB) ExpenseRepoInterface {
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

func (r *ExpenseRepo) Get(offset, limit int, f ExpenseFilter) ([]model.Expense, error) {
	var expenses []model.Expense

rows, err := r.DB.Query(`
	SELECT id, date, amount, note
	FROM expenses
	WHERE ($1 = 0 OR amount >= $1)
	AND ($2 = 0 OR amount <= $2)
	AND ($3 = '0001-01-01'::date OR date >= $3::date)
	AND ($4 = '0001-01-01'::date OR date <= $4::date)
	ORDER BY id DESC
	OFFSET $5
	LIMIT $6
	`, f.Min, f.Max, f.From, f.To, offset, limit,
)
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

func (r *ExpenseRepo) Count(f ExpenseFilter) (int, error) {
	var total int

err := r.DB.QueryRow(`
	SELECT COUNT(*)
	FROM expenses
	WHERE ($1 = 0 OR amount >= $1)
	AND ($2 = 0 OR amount <= $2)
	AND ($3 = '0001-01-01'::date OR date >= $3::date)
	AND ($4 = '0001-01-01'::date OR date <= $4::date)
	`, f.Min, f.Max, f.From, f.To).Scan(&total)
	
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *ExpenseRepo) GetById(id int64) (*model.Expense, error) {
	var expense model.Expense
	err := r.DB.QueryRow(`
	SELECT id, date, amount, note
	FROM expenses
	WHERE id=$1`, id).Scan(
		&expense.ID, &expense.Date, &expense.Amount, &expense.Note)

	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
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
		return nil, errors.ErrNotFound
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
		return nil, errors.ErrNotFound
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
		return errors.ErrNotFound
	}
	if err != nil {
		return err
	}

	return nil
}
