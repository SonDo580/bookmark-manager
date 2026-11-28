package db

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
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

type errorMatcher func(err error) bool

func getErrorFilter(matcher errorMatcher, targetErr error) ErrorFilter {
	return func(err error) (bool, error) {
		if err == nil {
			return false, nil
		}
		if matcher(err) {
			return true, targetErr
		}
		return false, nil
	}
}

// (The app uses Postgres while the tests use SQLite)

func UniqueViolationFilter(pgConstraint, sqliteTarget string, targetErr error) ErrorFilter {
	matcher := func(err error) bool {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == pgConstraint {
			return true
		}
		return strings.Contains(err.Error(), "UNIQUE constraint failed: "+sqliteTarget)
	}
	return getErrorFilter(matcher, targetErr)
}

func NoRecordFilter(targetErr error) ErrorFilter {
	matcher := func(err error) bool {
		return errors.Is(err, gorm.ErrRecordNotFound)
	}
	return getErrorFilter(matcher, targetErr)
}
