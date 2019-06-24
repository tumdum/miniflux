package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"miniflux.app/database"
	"miniflux.app/model"
)

const (
	testUser     = "foo"
	testPassword = "bar"
)

func MustCreateInMemoryStorage() *Storage {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	if err := database.Migrate(db); err != nil {
		panic(err)
	}
	return NewStorage(db)
}

func TestNoUserExistsInEmptyStorage(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	if storage.UserExists(testUser) {
		t.Fatalf("User '%v' should not exist", testUser)
	}
}

func TestAfterCreatingUserItExists(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user := model.User{
		Username: testUser,
		Password: testPassword,
	}
	if err := storage.CreateUser(&user); err != nil {
		t.Fatalf("Failed to create valid user: %v", err)
	}
	if user.ID == 0 {
		t.Fatalf("Failed to assign valid user ID")
	}
}
