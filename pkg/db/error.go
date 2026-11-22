package db

import (
	"errors"

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

func pgErrorFilter(errCode, constraintName string, targetErr error) ErrorFilter {
	return func(err error) (bool, error) {
		if err == nil {
			return false, nil
		}

		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return false, nil
		}

		if pgErr.Code != errCode {
			return false, nil
		}

		if pgErr.ConstraintName != constraintName {
			return false, nil
		}

		return true, targetErr
	}
}

func PGUniqueViolationFilter(constraintName string, targetErr error) ErrorFilter {
	return pgErrorFilter("23505", constraintName, targetErr)
}
