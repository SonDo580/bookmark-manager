package db

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateTestDB(t *testing.T) *gorm.DB {
	dsn := ":memory:"
	db, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		t.Fatalf("unable to connect to database: %v", err)
	}

	return db
}
