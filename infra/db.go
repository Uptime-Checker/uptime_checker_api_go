package infra

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
)

var DB *sql.DB

func ConnectDatabase() error {
	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.App.DatabaseHost, config.App.DatabasePort, config.App.DatabaseUser, config.App.DatabasePassword,
		config.App.DatabaseSchema)
	// Get a database handle.
	var err error
	DB, err = sql.Open("postgres", connectString)
	if err != nil {
		return err
	}

	return DB.Ping()
}
