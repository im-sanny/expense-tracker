package repository

import (
	"database/sql"
	"expense-tracker/internal/model"
	"expense-tracker/pkg/errors"
	"fmt"
	"strings"
	"time"
)

type ExpenseFilter struct {
	Min    int       `json:"min"`
	Max    int       `json:"max"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
	Search string    `json:"search"`
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

func (r *ExpenseRepo) buildWhereClause(f ExpenseFilter) (string, []interface{}) {
	var (
		whereParts []string
		args       []interface{}
	)
	argCount := 1

	// min amount filter
	if f.Min > 0 {
		whereParts = append(whereParts, fmt.Sprintf("amount >= $%d", argCount))
		args = append(args, f.Min)
		argCount++
	}

	// max amount filter
	if f.Max > 0 {
		whereParts = append(whereParts, fmt.Sprintf("amount <= $%d", argCount))
		args = append(args, f.Max)
		argCount++
	}

	// from data filter
	if !f.From.IsZero() && f.From != time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) {
		whereParts = append(whereParts, fmt.Sprintf("date >= $%d", argCount))
		args = append(args, f.From)
		argCount++
	}

	// to data filter
	if !f.To.IsZero() && f.To != time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) {
		whereParts = append(whereParts, fmt.Sprintf("date <= $%d", argCount))
		args = append(args, f.To)
		argCount++
	}

	// search (note) filter
	if strings.TrimSpace(f.Search) != "" {
		whereParts = append(whereParts, fmt.Sprintf("note ILIKE $%d", argCount))
		args = append(args, "%"+strings.TrimSpace(f.Search)+"%")
		argCount++
	}

	// combine WHERE parts
	if len(whereParts) > 0 {
		return " WHERE " + strings.Join(whereParts, " AND "), args
	}
	return "", args
}

func (r *ExpenseRepo) Get(offset, limit int, f ExpenseFilter) ([]model.Expense, error) {
	// build base query
	baseQuery := "SELECT id, date, amount, note FROM expenses"

	// build dynamic where clause
	whereClause, whereArgs := r.buildWhereClause(f)

	// add ORDER, OFFSET, LIMIT
	query := baseQuery + whereClause + " ORDER BY id DESC OFFSET $%d LIMIT $%d"

	// append pagination args
	args := append(whereArgs, offset, limit)

	// format the final query with correct placeholders
	query = fmt.Sprintf(query, len(whereArgs)+1, len(whereArgs)+2)

	// execute query
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// scan results
	var expenses []model.Expense
	for rows.Next() {
		var expense model.Expense
		if err := rows.Scan(&expense.ID, &expense.Date, &expense.Amount, &expense.Note); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, rows.Err()
}

func (r *ExpenseRepo) Count(f ExpenseFilter) (int, error) {
	// build base query
	baseQuery := "SELECT COUNT(*) FROM expenses"

	// build dynamic WHERE clause
	whereClause, whereArgs := r.buildWhereClause(f)

	// combine
	query := baseQuery + whereClause

	// execute query
	var total int
	err := r.DB.QueryRow(query, whereArgs...).Scan(&total)
	return total, err
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
