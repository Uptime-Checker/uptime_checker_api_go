package watchdog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/imroc/req/v3"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

const (
	appName     = "uptime_checker"
	website     = "https://uptimecheckr.com"
	version     = "v1.0"
	maxRedirect = 5
)

type HitErr struct {
	Type resource.ErrorLogType
	Text string
}

type HitResponse struct {
	StatusCode  int
	Duration    *int32
	Size        *int32
	ContentType *string
	Body        *string
	Headers     *map[string]string
	Traces      req.TraceInfo
}

func (w *WatchDog) Hit(
	ctx context.Context,
	uri, method string, body, username, password *string,
	bodyFormat resource.MonitorBodyFormat,
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

	requestContentType := w.monitorService.GetRequestContentType(bodyFormat, headers)
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

	lgr.Print(tracingID, "hitting =>", method, uri, "timeout", fmt.Sprintf("%ds", timeout))
	resp, err := request.Send(method, uri)
	if err != nil {
		return nil, w.getError(err)
	}
	// Return status code, response body, response headers, response size, trace info
	return w.getResponse(resp)
}

func (w *WatchDog) getError(err error) *HitErr {
	errorText := err.Error()
	errorLogType := resource.ErrorLogTypeUnknown

	var netError net.Error
	var urlError *url.Error
	var addrError *net.AddrError
	var dnsError *net.DNSError
	var invalidAddrError *net.InvalidAddrError
	var parseError *net.ParseError
	var unknownNetworkError *net.UnknownNetworkError
	var syscallError *os.SyscallError

	if errors.As(err, &netError) && netError.Timeout() {
		errorText = constant.ErrRemoteServerMaxTimeout
		errorLogType = resource.ErrorLogTypeTimeout
	} else {
		clientError := errors.Unwrap(err)
		if clientError != nil {
			errorText = clientError.Error()
			serverError := errors.Unwrap(clientError)
			if serverError != nil {
				errorText = serverError.Error()
			}
		}

		switch {
		case errors.As(err, &addrError):
			errorLogType = resource.ErrorLogTypeAddr
		case errors.As(err, &dnsError):
			errorLogType = resource.ErrorLogTypeDNS
		case errors.As(err, &invalidAddrError):
			errorLogType = resource.ErrorLogTypeInvalidAddr
		case errors.As(err, &parseError):
			errorLogType = resource.ErrorLogTypeParse
		case errors.As(err, &unknownNetworkError):
			errorLogType = resource.ErrorLogTypeUnknownNetwork
		case errors.As(err, &syscallError):
			errorLogType = resource.ErrorLogTypeSyscall
		case errors.As(err, &urlError):
			errorLogType = resource.ErrorLogTypeURL
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

	size := int32(resp.ContentLength)
	duration := int32(resp.TotalTime().Milliseconds())
	hitResponse = &HitResponse{
		StatusCode:  resp.GetStatusCode(),
		Duration:    &duration,
		Size:        &size,
		ContentType: pkg.StringPointer(resp.GetContentType()),
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
	if pkg.IsEmpty(*hitResponse.ContentType) {
		hitResponse.ContentType = pkg.StringPointer(http.DetectContentType(respBody))
	}

	return hitResponse, hitErr
}

func (w *WatchDog) getResponseHeaders(headers http.Header) *map[string]string {
	respHeader := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			respHeader[key] = values[0]
		}
	}
	return &respHeader
}
