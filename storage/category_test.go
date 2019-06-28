package storage

import (
	"testing"

	"miniflux.app/model"
)

func TestEmptyStorageHasNoCategories(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	if storage.CategoryExists(1, 1) {
		t.Fatalf("Unexpected category exists")
	}
}

func TestAfterUserIsCreatedAllCategoryIsAlsoCreatedForIt(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user := model.User{
		Username: testUser,
		Password: testPassword,
	}
	noErr(t, storage.CreateUser(&user))
	cat, err := storage.CategoryByTitle(user.ID, "All")
	noErr(t, err)
	if !storage.CategoryExists(user.ID, cat.ID) {
		t.Fatalf("'All' category doesn't exist")
	}
}
