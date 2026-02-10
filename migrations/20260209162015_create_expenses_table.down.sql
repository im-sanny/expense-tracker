-- +migrate Down
DROP INDEX IF EXISTS expenses_date_idx;
DROP TABLE IF EXISTS expenses;
