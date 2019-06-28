package storage

import (
	"reflect"
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
	user := model.User{Username: testUser}
	noErr(t, storage.CreateUser(&user))
	cat, err := storage.CategoryByTitle(user.ID, "All")
	noErr(t, err)
	if !storage.CategoryExists(user.ID, cat.ID) {
		t.Fatalf("'All' category doesn't exist")
	}
}

func TestCreatingManyCategoriesShouldMakeThemVisible(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user1 := model.User{Username: testUser}
	user2 := model.User{Username: testUser + "!"}
	noErr(t, storage.CreateUser(&user1))
	noErr(t, storage.CreateUser(&user2))
	cat1 := model.Category{
		Title:     "cat1",
		UserID:    user1.ID,
		FeedCount: 5,
	}
	cat2 := model.Category{
		Title:     "cat2",
		UserID:    user2.ID,
		FeedCount: 335,
	}
	cat3 := model.Category{
		Title:     "cat3",
		UserID:    user1.ID,
		FeedCount: 15,
	}
	noErr(t, storage.CreateCategory(&cat1))
	noErr(t, storage.CreateCategory(&cat2))
	noErr(t, storage.CreateCategory(&cat3))
	allForU1, err := storage.Categories(user1.ID)
	noErr(t, err)
	allForU2, err := storage.Categories(user2.ID)
	noErr(t, err)
	if len(allForU1) != 3 || len(allForU2) != 2 {
		t.Fatalf("Unexpected lengths of '%v' and '%v'", allForU1, allForU2)
	}
	for _, name := range []string{"All", "cat1", "cat3"} {
		found := false
		for _, cat := range allForU1 {
			if cat.Title == name {
				found = true
			}
		}
		if !found {
			t.Fatalf("Missing category '%v' for user '%v'", name, user1.Username)
		}
	}
	for _, name := range []string{"All", "cat2"} {
		found := false
		for _, cat := range allForU2 {
			if cat.Title == name {
				found = true
			}
		}
		if !found {
			t.Fatalf("Missing category '%v' for user '%v'", name, user1.Username)
		}
	}
}

func TestCategoryGetter(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user := model.User{Username: testUser}
	noErr(t, storage.CreateUser(&user))
	cat := model.Category{
		Title:     "cat",
		UserID:    user.ID,
		FeedCount: 0,
	}
	err := storage.CreateCategory(&cat)
	noErr(t, err)
	got, err := storage.Category(user.ID, cat.ID)
	noErr(t, err)
	if !reflect.DeepEqual(cat, *got) {
		t.Fatalf("Expected '%v' got '%v'", &cat, got)
	}
}

func TestRemoveCategory(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user := model.User{Username: testUser}
	noErr(t, storage.CreateUser(&user))
	cat, err := storage.CategoryByTitle(user.ID, "All")
	noErr(t, err)
	noErr(t, storage.RemoveCategory(user.ID, cat.ID))
	if storage.RemoveCategory(user.ID, cat.ID) == nil {
		t.Fatalf("Removing not existing category didn't fail")
	}
}

func TestCategoriesWithFeedCount(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user1 := model.User{Username: testUser}
	noErr(t, storage.CreateUser(&user1))
	cat1 := model.Category{
		Title:     "cat1",
		UserID:    user1.ID,
		FeedCount: 5,
	}
	cat2 := model.Category{
		Title:     "cat2",
		UserID:    user1.ID,
		FeedCount: 335,
	}
	cat3 := model.Category{
		Title:     "cat3",
		UserID:    user1.ID,
		FeedCount: 15,
	}
	noErr(t, storage.CreateCategory(&cat1))
	noErr(t, storage.CreateCategory(&cat2))
	noErr(t, storage.CreateCategory(&cat3))
	all, err := storage.CategoriesWithFeedCount(user1.ID)
	noErr(t, err)
	if len(all) != 4 {
		t.Fatal()
	}
	expected := map[string]int{
		"All":  0,
		"cat1": 5,
		"cat2": 335,
		"cat3": 15,
	}
	got := map[string]int{}
	for _, cat := range all {
		got[cat.Title] = cat.FeedCount
	}
	if reflect.DeepEqual(expected, got) {
		t.Fatalf("Expected '%v', got '%v'", expected, got)
	}
}

func TestFirstCategory(t *testing.T) {
	storage := MustCreateInMemoryStorage()
	defer storage.Close()
	user1 := model.User{Username: testUser}
	noErr(t, storage.CreateUser(&user1))
	cat1 := model.Category{
		Title:     "a",
		UserID:    user1.ID,
		FeedCount: 5,
	}
	cat2 := model.Category{
		Title:     "b",
		UserID:    user1.ID,
		FeedCount: 335,
	}
	cat3 := model.Category{
		Title:     "c",
		UserID:    user1.ID,
		FeedCount: 15,
	}
	noErr(t, storage.CreateCategory(&cat3))
	noErr(t, storage.CreateCategory(&cat1))
	noErr(t, storage.CreateCategory(&cat2))
	cat, err := storage.CategoryByTitle(user1.ID, "All")
	noErr(t, err)
	noErr(t, storage.RemoveCategory(user1.ID, cat.ID))
	first, err := storage.FirstCategory(user1.ID)
	noErr(t, err)
	if first.Title != "a" {
		t.Fatal()
	}
}