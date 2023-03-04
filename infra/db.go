package infra

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
)

var DB *sql.DB

func ConnectDatabase(enableLogging bool) error {
	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.App.DatabaseHost, config.App.DatabasePort, config.App.DatabaseUser, config.App.DatabasePassword,
		config.App.DatabaseSchema)
	// Get a database handle.
	var err error
	DB, err = sql.Open("postgres", connectString)
	if err != nil {
		return err
	}

	if enableLogging {
		postgres.SetQueryLogger(func(ctx context.Context, queryInfo postgres.QueryInfo) {
			// Depending on how the statement is executed, RowsProcessed is:
			//   - Number of rows returned for Query() and QueryContext() methods
			//   - RowsAffected() for Exec() and ExecContext() methods
			//   - Always 0 for Rows() method.
			fmt.Printf("|>----------------------------------------------------- %s|> processed [%d] - in %.2fs\n",
				queryInfo.Statement.DebugSql(), queryInfo.RowsProcessed, queryInfo.Duration.Seconds())
		})
	}

	return DB.Ping()
}

func StartTransaction(ctx context.Context) (*sql.Tx, error) {
	return DB.BeginTx(ctx, nil)
}

func CommitTransaction(transaction *sql.Tx) error {
	return transaction.Commit()
}
