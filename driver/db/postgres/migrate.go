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

func TestMigrate(db *sql.DB) {
	cwd, _ := os.Getwd()
	migrationDir := &migrate.FileMigrationSource{
		Dir: path.Join(cwd, "../../driver/db/postgres/migrations"),
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

func Seed(db *sql.DB) {
	db.Exec(`INSERT INTO users(id, rate_limit, quota, created_at, updated_at)
VALUES('123456', 2, 10, NOW(), NOW()) ON CONFLICT (id) DO UPDATE SET rate_limit = 2;`)

	db.Exec(`INSERT INTO user_usage(user_id, quota, quota_usage,
                       start_date, end_date, created_at, updated_at) 
VALUES('123456', 10, 0, NOW(), NOW() + interval '1 month', NOW(), NOW()) ON CONFLICT (user_id) DO UPDATE SET quota_usage = 0;`)
}
