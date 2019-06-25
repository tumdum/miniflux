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
	if !storage.UserExists(testUser) {
		t.Fatalf("Created user '%v' should exist", testUser)
	}
}

func TestAfterCreatingManyUsersTheyAllExists(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	userNames := []string{"user1", "user2", "user3"}
	ids := map[int64]struct{}{}
	for _, userName := range userNames {
		user := model.User{
			Username: userName,
			Password: testPassword,
		}
		if err := storage.CreateUser(&user); err != nil {
			t.Fatalf("Failed to create valid user: %v", err)
		}
		if user.ID == 0 {
			t.Fatalf("Failed to assign valid user ID")
		}
		ids[user.ID] = struct{}{}
		if !storage.UserExists(userName) {
			t.Fatalf("Created user '%v' should exist", userName)
		}
	}
	if len(userNames) != len(ids) {
		t.Fatalf("Expected %d unique ids, got %v", len(userNames), ids)
	}
}

func TestRemovingExistingUser(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user := model.User{
		Username: testUser,
		Password: testPassword,
	}
	if err := storage.CreateUser(&user); err != nil {
		t.Fatalf("Failed to create valid user: %v", err)
	}
	if err := storage.RemoveUser(user.ID); err != nil {
		t.Fatalf("Failed to remove valid user: %v", err)
	}
	if storage.UserExists(testUser) {
		t.Fatalf("User '%v' shouldn't exist", testUser)
	}
}

func TestRemovingNotExistingUserFails(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	if err := storage.RemoveUser(1); err == nil {
		t.Fatalf("Romeving not existing users didn't fail")
	}
}
