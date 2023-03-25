package watchdog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
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
	Text string
}

type HitResponse struct {
	StatusCode  int
	Body        *string
	Headers     *map[string]string
	Size        int64
	ContentType *string
	Traces      *string
}

func Hit(
	ctx context.Context,
	url, method string, body, username, password *string,
	bodyFormat *resource.MonitorBodyFormat,
	headers *map[string]string,
	timeout int,
	followRedirect bool,
) (*HitResponse, *HitErr) {
	tracingID := pkg.GetTracingID(ctx)

	agent := fmt.Sprintf("%s_agent/%s (%s)", appName, version, website)
	client := req.C().SetTimeout(time.Duration(timeout) * time.Second).SetUserAgent(agent)
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

	lgr.Default.Print(tracingID, "Hitting =>", method, url, "timeout", timeout, "s")
	resp, err := request.Send(method, url)
	if err != nil {
		return nil, getError(err)
	}
	// Return status code, response body, response headers, response size, trace info
	return getResponse(resp)
}

func getRequestContentType(bodyFormat *resource.MonitorBodyFormat, headers *map[string]string) string {
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

func getError(err error) *HitErr {
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

func getResponse(resp *req.Response) (*HitResponse, *HitErr) {
	var respTrace *string

	var hitErr *HitErr
	var hitResponse *HitResponse

	if resp.Body == nil {
		return hitResponse, hitErr
	}

	traceInfo, err := json.Marshal(resp.TraceInfo())
	if err == nil {
		stringTraceInfo := string(traceInfo)
		respTrace = &stringTraceInfo
	}

	headers := getResponseHeaders(resp.Header)
	contentType := resp.GetContentType()
	hitResponse = &HitResponse{
		StatusCode:  resp.GetStatusCode(),
		Size:        resp.ContentLength,
		Traces:      respTrace,
		Headers:     headers,
		ContentType: &contentType,
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		hitErr = &HitErr{
			Type: resource.ErrorLogTypeResponseMalformed,
			Text: constant.ErrResponseMalformed,
		}
		return hitResponse, hitErr
	} else {
		stringBody := string(respBody)
		hitResponse.Body = &stringBody

		if hitResponse.ContentType == nil {
			mimeType := http.DetectContentType(respBody)
			hitResponse.ContentType = &mimeType
		}
	}

	return hitResponse, hitErr
}

func getResponseHeaders(headers http.Header) *map[string]string {
	respHeader := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			respHeader[key] = values[0]
		}
	}
	return &respHeader
}
