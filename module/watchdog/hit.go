package watchdog

import (
	"context"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
)

func Hit(
	ctx context.Context,
	url, method, body *string,
	bodyFormat *resource.MonitorBodyFormat,
	headers *map[string]string,
	timeout int,
	followRedirect bool,
) {

}
