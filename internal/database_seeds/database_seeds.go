package database_seeds

import (
	"errors"
	"fmt"
	"github.com/BudjeeApp/budget-core/internal/config"
	"os"
)

// Refresh truncates, then loads all tables with data
func Refresh() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "DEVELOPMENT" {
		config.Logger.Fatal(
			errors.New(
				"seedDatabase is available only in DEVELOPMENT. " +
					"Make sure you have APP_ENV set to DEVELOPMENT",
			).Error(),
		)
	}

	db := config.GetDB()
	tables := []string{
		"transactions",
		"categories",
		"accounts",
	}
	for _, t := range tables {
		statement := fmt.Sprintf("TRUNCATE %v CASCADE", t)
		_, err := db.Exec(statement)
		if err != nil {
			config.Logger.Fatal(err.Error())
		}
	}
}
