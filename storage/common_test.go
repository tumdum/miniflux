package storage

import (
	"database/sql"
	"testing"

	"miniflux.app/database"
)

func noErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func MustCreateStorage(path string) *Storage {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	if err := database.Migrate(db); err != nil {
		panic(err)
	}
	return NewStorage(db)
}

func MustCreateInMemoryStorage() *Storage {
	return MustCreateStorage(":memory:")
}
