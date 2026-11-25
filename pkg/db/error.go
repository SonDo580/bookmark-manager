package db

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorFilter func(err error) (matched bool, mappedErr error)

func CatchDBError(filters []ErrorFilter, err error) error {
	for _, filter := range filters {
		matched, mappedErr := filter(err)
		if matched {
			return mappedErr
		}
	}

	return err
}

// (The app uses Postgres while the tests use SQLite)

func UniqueViolationFilter(pgConstraint, sqliteTarget string, targetErr error) ErrorFilter {
	return func(err error) (bool, error) {
		if err == nil {
			return false, nil
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == pgConstraint {
			return true, targetErr
		}

		if strings.Contains(err.Error(), "UNIQUE constraint failed: "+sqliteTarget) {
			return true, targetErr
		}

		return false, nil
	}
}
