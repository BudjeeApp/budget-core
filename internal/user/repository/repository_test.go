package repository_test

import (
	"github.com/BudjeeApp/budget-core/internal/test"
	. "github.com/BudjeeApp/budget-core/internal/user/repository"
	_ "github.com/lib/pq"
	"testing"
)

var (
	repo Repository
)

func setupTest(t *testing.T) {
	t.Helper()
	testRepo, err := test.GetDBInstance()
	if err != nil {
		t.Fatal(err.Error())
	}
	repo = testRepo.User
	test.RemoveUsers(t)
}

func teardownTest(t *testing.T) {
	test.RemoveUsers(t)
}
