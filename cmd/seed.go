package cmd

import (
	"github.com/BudjeeApp/budget-core/internal/database_seeds"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(seedDatabaseCmd)
}

var seedDatabaseCmd = &cobra.Command{
	Use: "seedDatabase",
	Run: func(cmd *cobra.Command, args []string) {
		CreateDatabaseSeeds()
	},
}

func CreateDatabaseSeeds() {
	database_seeds.Refresh()
}
