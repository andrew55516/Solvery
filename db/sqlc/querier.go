// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"
)

type Querier interface {
	AddUserCredit(ctx context.Context, arg AddUserCreditParams) (User, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, email string) (User, error)
	ListAllEntries(ctx context.Context, arg ListAllEntriesParams) ([]Entry, error)
	ListUserEntries(ctx context.Context, arg ListUserEntriesParams) ([]Entry, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
}

var _ Querier = (*Queries)(nil)