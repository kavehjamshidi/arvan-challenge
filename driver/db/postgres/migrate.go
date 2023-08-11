package postgres

import (
	"database/sql"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

func Migrate(db *sql.DB) {
	cwd, _ := os.Getwd()
	migrationDir := &migrate.FileMigrationSource{
		Dir: path.Join(cwd, "/driver/db/postgres/migrations"),
	}
	executor := &migrate.MigrationSet{
		TableName: viper.GetString("DB_MIGRATION_TABLE"),
	}
	n, err := executor.Exec(db, "postgres", migrationDir, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("migration encountered a problem: %v\n", err))
	}
	log.Printf("applied %d migrations!\n", n)
}
