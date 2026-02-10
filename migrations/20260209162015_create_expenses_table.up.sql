-- +migrate Up
CREATE TABLE IF NOT EXISTS expenses (
    id      BIGSERIAL PRIMARY KEY,
    date    DATE      NOT NULL,
    amount  INTEGER   NOT NULL CHECK (amount >= 0),
    note    TEXT
);

CREATE INDEX IF NOT EXISTS expenses_date_idx ON expenses (date);
