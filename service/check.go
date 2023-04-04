package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/imroc/req/v3"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type CheckService struct {
	checkDomain *domain.CheckDomain
}

func NewCheckService(checkDomain *domain.CheckDomain) *CheckService {
	return &CheckService{checkDomain: checkDomain}
}

func (c *CheckService) Update(
	ctx context.Context,
	tx *sql.Tx,
	check *model.Check,
	success bool,
	duration, size *int32,
	contentType *string,
	body *string,
	headers *map[string]string,
	traces req.TraceInfo,
) (*model.Check, error) {
	if headers != nil && len(*headers) > 0 {
		jsonHeaders, err := json.Marshal(*headers)
		if err != nil {
			return nil, err
		}
		check.Headers = pkg.StringPointer(string(jsonHeaders))
	}

	check.Body = body
	check.Success = success
	check.Duration = duration
	check.ContentSize = size
	check.ContentType = contentType

	traceInfo := make(map[string]string)

	traceInfo["TotalTime"] = traces.TotalTime.String()
	traceInfo["DNSLookupTime"] = traces.TotalTime.String()
	traceInfo["TCPConnectTime"] = traces.TotalTime.String()
	traceInfo["TLSHandshakeTime"] = traces.TotalTime.String()
	traceInfo["FirstResponseTime"] = traces.TotalTime.String()
	traceInfo["ResponseTime"] = traces.TotalTime.String()

	jsonTraces, err := json.Marshal(traceInfo)
	if err != nil {
		return nil, err
	}
	check.Traces = pkg.StringPointer(string(jsonTraces))

	return c.checkDomain.Update(ctx, tx, check)
}
