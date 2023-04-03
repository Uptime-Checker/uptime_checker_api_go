package infra

import (
	"context"
	"database/sql"
	"fmt"

	// pgx is the postgres driver needed to be imported
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/go-jet/jet/v2/postgres"

	"github.com/fatih/color"
	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

var DB *sql.DB

func ConnectDatabase(ctx context.Context, enableLogging bool) error {
	tracingID := pkg.GetTracingID(ctx)
	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.App.DatabaseHost, config.App.DatabasePort, config.App.DatabaseUser, config.App.DatabasePassword,
		config.App.DatabaseSchema)
	// Get a database handle.
	var err error
	DB, err = sql.Open("pgx", connectString)
	if err != nil {
		return err
	}

	postgres.SetQueryLogger(func(ctx context.Context, queryInfo postgres.QueryInfo) {
		// Depending on how the statement is executed, RowsProcessed is:
		//   - Number of rows returned for Query() and QueryContext() methods
		//   - RowsAffected() for Exec() and ExecContext() methods
		//   - Always 0 for Rows() method.
		if enableLogging {
			c0 := color.New(color.FgHiGreen)
			_, _ = c0.Print("|>-----------------------------------------------------")
			c1 := color.New(color.FgCyan)
			_, _ = c1.Printf("%s", queryInfo.Statement.DebugSql())
			c2 := color.New(color.FgHiRed).Add(color.Underline)
			_, _ = c2.Printf("|> processed [%d] - in %.2fs\n", queryInfo.RowsProcessed, queryInfo.Duration.Seconds())
		}
		if queryInfo.Duration.Seconds() > 1 {
			lgr.Warn(tracingID, "slow query, time:", fmt.Sprintf("%d ms", queryInfo.Duration.Milliseconds()),
				"rows processed", queryInfo.RowsProcessed,
				"sql: ", queryInfo.Statement.DebugSql())
		}
	})

	return DB.PingContext(ctx)
}

func StartTransaction(ctx context.Context) (*sql.Tx, error) {
	tracingID := pkg.GetTracingID(ctx)
	lgr.Print(tracingID, "start transaction")

	return DB.BeginTx(ctx, nil)
}

func CommitTransaction(ctx context.Context, transaction *sql.Tx) error {
	tracingID := pkg.GetTracingID(ctx)
	lgr.Print(tracingID, "commit transaction")

	return transaction.Commit()
}

func RollbackTransaction(ctx context.Context, transaction *sql.Tx) error {
	tracingID := pkg.GetTracingID(ctx)
	lgr.Print(tracingID, "rollback transaction")

	return transaction.Rollback()
}

// Transaction creates a transaction and calls f.
// When it is finished, it cleans up the transaction. If an error occurred it
// attempts to rollback, if not it commits.
// Queries inside the transaction cannot run concurrently
// https://github.com/jackc/pgx/issues/1256
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
