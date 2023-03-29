package watchdog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
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
	Text string
}

type HitResponse struct {
	StatusCode  int
	Duration    int64
	Size        int64
	ContentType string
	Body        *string
	Headers     map[string]string
	Traces      req.TraceInfo
}

func (w *WatchDog) Hit(
	ctx context.Context,
	url, method string, body, username, password *string,
	bodyFormat *resource.MonitorBodyFormat,
	headers *map[string]string,
	timeout int32,
	followRedirect bool,
) (*HitResponse, *HitErr) {
	tracingID := pkg.GetTracingID(ctx)

	agent := fmt.Sprintf("%s_agent/%s (%s)", appName, version, website)
	client := req.C().SetTimeout(time.Duration(timeout) * time.Second).SetUserAgent(agent).DisableAutoReadResponse()
	if followRedirect {
		client.SetRedirectPolicy(req.MaxRedirectPolicy(maxRedirect))
	}

	request := client.R().SetContext(ctx).EnableTrace()

	requestContentType := getRequestContentType(bodyFormat, headers)
	if requestContentType != "" {
		request.SetContentType(requestContentType)
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

	lgr.Print(tracingID, "hitting =>", method, url, "timeout", fmt.Sprintf("%ds", timeout))
	resp, err := request.Send(method, url)
	if err != nil {
		return nil, w.getError(err)
	}
	// Return status code, response body, response headers, response size, trace info
	return w.getResponse(resp)
}

func getRequestContentType(bodyFormat *resource.MonitorBodyFormat, headers *map[string]string) string {
	contentType := resource.MonitorBodyFormatNoBody.String()
	if bodyFormat != nil {
		contentType = bodyFormat.String()
	}

	if headers != nil && len(*headers) > 0 {
		for key, value := range *headers {
			if strings.EqualFold(key, contentTypeHeaderKey) {
				contentType = value
			}
		}
	}

	return contentType
}

func (w *WatchDog) getError(err error) *HitErr {
	clientError := errors.Unwrap(err)

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

	return &HitErr{Type: errorLogType, Text: errorText}
}

func (w *WatchDog) getResponse(resp *req.Response) (*HitResponse, *HitErr) {
	var hitErr *HitErr
	var hitResponse *HitResponse

	if resp.Body == nil {
		return hitResponse, hitErr
	}

	hitResponse = &HitResponse{
		StatusCode:  resp.GetStatusCode(),
		Duration:    resp.TotalTime().Milliseconds(),
		Size:        resp.ContentLength,
		ContentType: resp.GetContentType(),
		Headers:     w.getResponseHeaders(resp.Header),
		Traces:      resp.TraceInfo(),
	}

	respBody, err := io.ReadAll(resp.Body)
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			sentry.CaptureException(err)
		}
	}(resp.Body)

	if err != nil {
		hitErr = &HitErr{Type: resource.ErrorLogTypeResponseMalformed, Text: constant.ErrResponseMalformed}
		return hitResponse, hitErr
	}
	stringBody := string(respBody)
	hitResponse.Body = &stringBody
	if hitResponse.ContentType == "" {
		hitResponse.ContentType = http.DetectContentType(respBody)
	}

	return hitResponse, hitErr
}

func (w *WatchDog) getResponseHeaders(headers http.Header) map[string]string {
	respHeader := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			respHeader[key] = values[0]
		}
	}
	return respHeader
}
