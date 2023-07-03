// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: entries.sql

package db

import (
	"context"
)

const createEntry = `-- name: CreateEntry :one
    INSERT INTO entries (user_email, amount, comment) VALUES ($1, $2, $3) RETURNING id, user_email, amount, comment, created_at
`

type CreateEntryParams struct {
	UserEmail string `json:"user_email"`
	Amount    int32  `json:"amount"`
	Comment   string `json:"comment"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.UserEmail, arg.Amount, arg.Comment)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.UserEmail,
		&i.Amount,
		&i.Comment,
		&i.CreatedAt,
	)
	return i, err
}

const listAllEntries = `-- name: ListAllEntries :many
    SELECT id, user_email, amount, comment, created_at FROM entries LIMIT $1 OFFSET $2
`

type ListAllEntriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAllEntries(ctx context.Context, arg ListAllEntriesParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listAllEntries, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.UserEmail,
			&i.Amount,
			&i.Comment,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUserEntries = `-- name: ListUserEntries :many
    SELECT id, user_email, amount, comment, created_at FROM entries WHERE user_email = $1 LIMIT $2 OFFSET $3
`

type ListUserEntriesParams struct {
	UserEmail string `json:"user_email"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

func (q *Queries) ListUserEntries(ctx context.Context, arg ListUserEntriesParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listUserEntries, arg.UserEmail, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.UserEmail,
			&i.Amount,
			&i.Comment,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}