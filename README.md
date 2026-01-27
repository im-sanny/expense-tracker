# Expense Tracker API

A simple RESTful expense tracker built with Go and PostgreSQL.
This project is mainly for learning core backend concepts like CRUD, REST semantics, and partial updates using PATCH.

## Features

- Create, read, update, and delete expenses
- Proper PUT vs PATCH behavior
- Partial updates using SQL `COALESCE`
- Clean and minimal HTTP handlers

## Tech Stack

- Go (net/http)
- PostgreSQL
- database/sql

## Expense Model

```json
{
  "id": 1,
  "date": "2026-01-26T18:30:00Z",
  "amount": 890,
  "note": "Coffee + croissant"
}
```
