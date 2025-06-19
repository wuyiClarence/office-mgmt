package migration

import (
	"errors"
	"fmt"
	"log"

	"platform-backend/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func AutoMigration() {
	m, err := migrate.New(
		"file://./migration",
		fmt.Sprintf("mysql://%s", config.MyConfig.MysqlConfig.Dsn),
	)
	if err != nil {
		panic(err)
	}

	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}

	log.Println("Database migration completed successfully.")
	return
}
