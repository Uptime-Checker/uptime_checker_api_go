package watchdog

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
)

const (
	appName = "uptime_checker"
	website = "https://uptimecheckr.com"
	version = "v1.0"
)

func Hit(
	ctx context.Context,
	url, method string, body, username, password *string,
	bodyFormat *resource.MonitorBodyFormat,
	headers *map[string]string,
	timeout int,
	followRedirect bool,
) {
	agent := fmt.Sprintf("%s_agent/%s (%s)", appName, version, website)
	client := req.C().SetTimeout(time.Duration(timeout) * time.Second).SetUserAgent(agent)
	client.DevMode()

	request := client.R()
	request.SetContext(ctx)

	if username != nil && password != nil {
		request.SetBasicAuth(*username, *password)
	}
	if headers != nil && len(*headers) > 0 {
		request.SetHeaders(*headers)
	}
	if body != nil {
		request.SetBody(*body)
	}

	_, err := request.Send(method, url)
	if err != nil {
		return
	}
}
