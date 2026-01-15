package cmd

import (
	"fmt"
	"log/slog"

	"github.com/msvens/mchess/internal/config"
	"github.com/msvens/mchess/internal/db"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long:  `Commands for managing the mchess database.`,
}

var dbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create database tables",
	Long:  `Run migrations to create all database tables.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		cfg.SetupLogger()

		slog.Info("Creating database tables...")
		if err := db.MigrateUp(cfg.DBConnectionString()); err != nil {
			slog.Error("Failed to create database", "error", err)
			return
		}
		slog.Info("Database tables created successfully")
	},
}

var dbDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete all database tables",
	Long:  `Drop all database tables. WARNING: This will delete all data!`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		cfg.SetupLogger()

		slog.Warn("Deleting all database tables...")
		if err := db.MigrateDownAll(cfg.DBConnectionString()); err != nil {
			slog.Error("Failed to delete database", "error", err)
			return
		}
		slog.Info("Database tables deleted")
	},
}

var dbUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade database schema",
	Long:  `Run any pending migrations to upgrade the database schema.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		cfg.SetupLogger()

		slog.Info("Upgrading database schema...")
		if err := db.MigrateUp(cfg.DBConnectionString()); err != nil {
			slog.Error("Failed to upgrade database", "error", err)
			return
		}
		slog.Info("Database schema upgraded successfully")
	},
}

var dbVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show database schema version",
	Long:  `Display the current database schema version.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		cfg.SetupLogger()

		version, dirty, err := db.GetVersion(cfg.DBConnectionString())
		if err != nil {
			slog.Error("Failed to get database version", "error", err)
			return
		}
		fmt.Printf("Database schema version: %d\n", version)
		if dirty {
			fmt.Println("  WARNING: Schema is in dirty state")
		}
	},
}

func init() {
	dbCmd.AddCommand(dbCreateCmd)
	dbCmd.AddCommand(dbDeleteCmd)
	dbCmd.AddCommand(dbUpgradeCmd)
	dbCmd.AddCommand(dbVersionCmd)
	rootCmd.AddCommand(dbCmd)
}
