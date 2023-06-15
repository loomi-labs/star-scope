/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/loomi-labs/star-scope/database"
	"github.com/spf13/cobra"
	"strings"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:     "database",
	Short:   "Database commands",
	Aliases: []string{"db"},
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Run: func(cmd *cobra.Command, args []string) {
		database.MigrateDb()
		database.InitDb()
	},
}

// createMigrationsCmd represents the createMigrations command
var createMigrationsCmd = &cobra.Command{
	Use:   "create-migrations",
	Short: "Create migrations based on ent/schema/*.go files",
	Long: `Create migrations based on ent/schema/*.go files

Example with custom db:
go run main.go createMigrations postgres://postgres:postgres@localhost:5432/star-scope-db?sslmode=disable&TimeZone=Europe/Zurich
`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			database.CreateMigrations(args[0])
		} else {
			database.CreateMigrations(strings.Replace(database.DbCon(), "5432", "5433", 1))
		}
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(migrateCmd)
	databaseCmd.AddCommand(createMigrationsCmd)
}
