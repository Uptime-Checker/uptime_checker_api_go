package infra

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/go-jet/jet/v2/postgres"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/log"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
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
	tracingID := pkg.GetTracingID(ctx)
	log.Default.Print(tracingID, "start transaction")

	return DB.BeginTx(ctx, nil)
}

func CommitTransaction(ctx context.Context, transaction *sql.Tx) error {
	tracingID := pkg.GetTracingID(ctx)
	log.Default.Print(tracingID, "commit transaction")

	return transaction.Commit()
}

func RollbackTransaction(ctx context.Context, transaction *sql.Tx) error {
	tracingID := pkg.GetTracingID(ctx)
	log.Default.Print(tracingID, "rollback transaction")

	return transaction.Rollback()
}

// Transaction creates a transaction and calls f.
// When it is finished, it cleans up the transaction. If an error occured it
// attempts to rollback, if not it commits.
func Transaction(ctx context.Context, f func(context.Context, *sql.Tx) error) error {
	tx, err := StartTransaction(ctx)
	if err != nil {
		return err
	}

	err = f(ctx, tx)
	if err != nil {
		rollbackError := RollbackTransaction(ctx, tx)
		if rollbackError != nil {
			sentry.CaptureException(rollbackError)
		}
		return err
	}
	return CommitTransaction(ctx, tx)
}
