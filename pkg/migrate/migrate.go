package migrate

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"

	"log"

	"github.com/PavelRadostev/train_trip/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration CLI",
	Long:  "Run, rollback, and inspect database schema migrations.",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running up migrations...")
		m := getMigrator()
		err := m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("❌ Migration up failed: %v", err)
		}

		version, dirty, _ := m.Version()
		fmt.Printf("✅ Migrated to version %d (dirty: %v)\n", version, dirty)
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback last migration",
	Run: func(cmd *cobra.Command, args []string) {
		m := getMigrator()
		err := m.Steps(-1)
		if err != nil {
			log.Fatalf("❌ Migration down failed: %v", err)
		}

		version, dirty, _ := m.Version()
		fmt.Printf("⬅️  Rolled back to version %d (dirty: %v)\n", version, dirty)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Run: func(cmd *cobra.Command, args []string) {
		m := getMigrator()
		version, dirty, err := m.Version()
		if err != nil {
			fmt.Printf("⚠️  Could not get version: %v\n", err)
		} else {
			fmt.Printf("📦 Current version: %d (dirty: %v)\n", version, dirty)
		}
	},
}

// Execute запускает корневую команду
func Execute() {
	fmt.Println("Starting database migration...")
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(versionCmd)
	_ = rootCmd.Execute()
}

func getMigrator() *migrate.Migrate {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		log.Fatalf("failed to create driver: %v", err)
	}

	// Формируем правильный file:// URL
	migrationsURL := "file://" + cfg.Migration.Dir
	m, err := migrate.NewWithDatabaseInstance(migrationsURL, "postgres", driver)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}

	return m
}
