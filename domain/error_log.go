package domain

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type ErrorLogDomain struct{}

func NewErrorLogDomain() *ErrorLogDomain {
	return &ErrorLogDomain{}
}

func (e *ErrorLogDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	errorLog *model.ErrorLog,
	errorLogType resource.ErrorLogType,
) (*model.ErrorLog, error) {
	if !errorLogType.Valid() {
		return nil, constant.ErrInvalidErrorLogType
	}
	errorLogTypeValue := int32(errorLogType)
	errorLog.Type = &errorLogTypeValue

	insertStmt := ErrorLog.INSERT(ErrorLog.MutableColumns.Except(ErrorLog.InsertedAt, ErrorLog.UpdatedAt)).
		MODEL(errorLog).
		RETURNING(ErrorLog.AllColumns)

	err := insertStmt.QueryContext(ctx, tx, errorLog)
	return errorLog, err
}
