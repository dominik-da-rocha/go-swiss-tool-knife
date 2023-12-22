package databox

import (
	"strconv"

	"github.com/dominik-da-rocha/go-toolbox/toolbox"
	"github.com/spf13/cobra"
)

var FlagDatabasePath string = "app.sqlite"
var FlagMigrationPath string = "migration"

func init() {
	MigrateCmd.AddCommand(downCmd)
	MigrateCmd.AddCommand(upCmd)
	MigrateCmd.AddCommand(toCmd)
	MigrateCmd.PersistentFlags().StringVar(&FlagDatabasePath, "database", FlagDatabasePath, "path to database")
	MigrateCmd.PersistentFlags().StringVar(&FlagMigrationPath, "migration", FlagMigrationPath, "path to migration files")
}

// MigrateCmd represents the migrate command
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates the database of the service",
	Long:  `Migrates the database of the service`,
	Run: func(cmd *cobra.Command, args []string) {
		d := OpenDb(FlagDatabasePath)
		defer d.Close()
		mig := NewMigration(d, FlagDatabasePath, FlagMigrationPath)
		GetMigrationVersion(mig)
	},
}

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "migrates down to previous version",
	Long:  `migrates down to previous version`,
	Run: func(cmd *cobra.Command, args []string) {
		d := OpenDb(FlagDatabasePath)
		defer d.Close()
		mig := NewMigration(d, FlagDatabasePath, FlagMigrationPath)
		GetMigrationVersion(mig)
		MigrateDown(mig)
	},
}

// downCmd represents the down command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "migrates up to the latest version",
	Long:  `migrates up to the latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		d := OpenDb(FlagDatabasePath)
		defer d.Close()
		mig := NewMigration(d, FlagDatabasePath, FlagMigrationPath)
		GetMigrationVersion(mig)
		MigrateUp(mig)
	},
}

// downCmd represents the down command
var toCmd = &cobra.Command{
	Use:   "to [version]",
	Short: "migrates to the version",
	Long:  `migrates to the version`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version64, err := strconv.ParseUint(args[0], 10, 32)
		toolbox.Uups(err)

		version := uint(version64)

		d := OpenDb(FlagDatabasePath)
		defer d.Close()
		mig := NewMigration(d, FlagDatabasePath, FlagMigrationPath)
		GetMigrationVersion(mig)
		MigrateTo(mig, version)
	},
}
