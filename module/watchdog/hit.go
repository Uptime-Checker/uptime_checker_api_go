package watchdog

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/imroc/req/v3"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

const (
	appName              = "uptime_checker"
	website              = "https://uptimecheckr.com"
	version              = "v1.0"
	maxRedirect          = 5
	contentTypeHeaderKey = "content-type"
)

type HitErr struct {
	Type resource.ErrorLogType
	Text *string
}

func Hit(
	ctx context.Context,
	url, method string, body, username, password *string,
	bodyFormat *resource.MonitorBodyFormat,
	headers *map[string]string,
	timeout int,
	followRedirect bool,
) {
	tracingID := pkg.GetTracingID(ctx)

	agent := fmt.Sprintf("%s_agent/%s (%s)", appName, version, website)
	client := req.C().SetTimeout(time.Duration(timeout) * time.Second).SetUserAgent(agent)
	if followRedirect {
		client.SetRedirectPolicy(req.MaxRedirectPolicy(maxRedirect))
	}

	var hitErr *HitErr

	client.OnAfterResponse(func(client *req.Client, resp *req.Response) error {
		if resp.Err != nil {
			hitErr = getError(resp)
			return nil
		}
		return nil
	})

	request := client.R().SetContext(ctx)

	contentType := getContentType(bodyFormat, headers)
	if contentType != "" {
		request.SetContentType(contentType)
	}
	if username != nil && password != nil {
		request.SetBasicAuth(*username, *password)
	}
	if headers != nil && len(*headers) > 0 {
		request.SetHeaders(*headers)
	}
	if body != nil {
		request.SetBody(*body)
	}

	lgr.Default.Print(tracingID, "Hitting =>", method, url, "timeout", timeout, "s")
	_, err := request.Send(method, url)
	if err != nil {
		return
	}

	// Return status code, response body, response headers, response size, trace info
}

func getContentType(bodyFormat *resource.MonitorBodyFormat, headers *map[string]string) string {
	contentType := resource.MonitorBodyFormatNoBody.String()
	if bodyFormat != nil {
		contentType = bodyFormat.String()
	}

	if headers != nil && len(*headers) > 0 {
		for key, value := range *headers {
			if strings.ToLower(key) == contentTypeHeaderKey {
				contentType = value
			}
		}
	}

	return contentType
}

func getError(resp *req.Response) *HitErr {
	if resp.Err == nil {
		return nil
	}

	clientError := errors.Unwrap(resp.Err)

	errorText := clientError.Error()
	errorLogType := resource.ErrorLogTypeUnknown

	if err, ok := clientError.(net.Error); ok && err.Timeout() {
		errorText = constant.ErrRemoteServerMaxTimeout
		errorLogType = resource.ErrorLogTypeTimeout
	} else if err != nil {
		// This was an error, but not a timeout
		serverError := errors.Unwrap(clientError)
		if serverError != nil {
			errorText = serverError.Error()

			if _, ok := serverError.(*net.AddrError); ok {
				errorLogType = resource.ErrorLogTypeAddr
			} else if _, ok := serverError.(*net.DNSError); ok {
				errorLogType = resource.ErrorLogTypeDNS
			} else if _, ok := serverError.(*net.InvalidAddrError); ok {
				errorLogType = resource.ErrorLogTypeInvalidAddr
			} else if _, ok := serverError.(*net.ParseError); ok {
				errorLogType = resource.ErrorLogTypeParse
			} else if _, ok := serverError.(*net.UnknownNetworkError); ok {
				errorLogType = resource.ErrorLogTypeUnknownNetwork
			} else if _, ok := serverError.(*os.SyscallError); ok {
				errorLogType = resource.ErrorLogTypeSyscall
			}
		}
	}

	return &HitErr{Type: errorLogType, Text: &errorText}
}
