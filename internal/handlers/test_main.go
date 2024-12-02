package handlers

import (
	"testing"
	"todo-api/internal/test"
)

func TestMain(m *testing.M) {
	// Set up test database
	db, err := test.SetupTestDB()
	if err != nil {
		panic(err)
	}

	// Run tests
	m.Run()

	// Clean up
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
