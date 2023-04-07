package service

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type ErrorLogService struct {
	errorLogDomain *domain.ErrorLogDomain
}

func NewErrorLogService(errorLogDomain *domain.ErrorLogDomain) *ErrorLogService {
	return &ErrorLogService{errorLogDomain: errorLogDomain}
}

func (e *ErrorLogService) Create(
	ctx context.Context,
	tx *sql.Tx, monitorID, checkID int64, assertionID *int64, text *string, errorLogType resource.ErrorLogType,
) (*model.ErrorLog, error) {
	errorLog := model.ErrorLog{
		MonitorID:   monitorID,
		CheckID:     checkID,
		AssertionID: assertionID,
		Text:        text,
	}

	return e.errorLogDomain.Create(ctx, tx, &errorLog, errorLogType)
}
