package storage

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"miniflux.app/database"
	"miniflux.app/model"
)

const (
	testUser     = "foo"
	testPassword = "bar"
)

func noErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

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
		noErr(t, storage.CreateUser(&user))
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
	noErr(t, storage.CreateUser(&user))
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

func TestGettingAllUsersFromEmptyStorageShouldReturnEmptySlice(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	if users, err := storage.Users(); len(users) > 0 {
		noErr(t, err)
		t.Fatalf("Expected no users, got %v", users)
	}
}

func TestGettingAllUsersShouldReturnThemAll(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	expected := map[string]struct{}{
		"abc": struct{}{},
		"def": struct{}{},
		"ghi": struct{}{},
	}
	for userName := range expected {
		user := model.User{
			Username: userName,
			Password: testPassword,
		}
		noErr(t, storage.CreateUser(&user))
	}
	allUsers, err := storage.Users()
	noErr(t, err)
	actual := map[string]struct{}{}
	for _, user := range allUsers {
		actual[user.Username] = struct{}{}
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected '%v', got '%v'", expected, actual)
	}
}

func TestUserBy(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user := model.User{
		Username: testUser,
		Password: testPassword,
	}
	noErr(t, storage.CreateUser(&user))

	got, err := storage.UserByUsername(testUser)
	noErr(t, err)
	if got.Username != testUser || got.ID != user.ID {
		t.Fatalf("Expected to get user '%v', got '%v'", testUser, got.Username)
	}
	got, err = storage.UserByID(user.ID)
	noErr(t, err)
	if got.Username != testUser || got.ID != user.ID {
		t.Fatalf("Expected to get user '%v', got '%v'", testUser, got.Username)
	}
}

func TestUserByFailsForNotExistingUser(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user, _ := storage.UserByUsername(testUser)
	if user != nil {
		t.Fatalf("Got not existing user: %#v", user)
	}
	user, _ = storage.UserByID(1)
	if user != nil {
		t.Fatalf("Got not existing user: %#v", user)
	}
}

func TestCheckPassword(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user := model.User{
		Username: testUser,
		Password: testPassword,
	}
	noErr(t, storage.CreateUser(&user))
	noErr(t, storage.CheckPassword(testUser, testPassword))
	if storage.CheckPassword(testUser, testPassword+"!") == nil {
		t.Fatalf("Check password didn't return error for incorrect password")
	}
}
