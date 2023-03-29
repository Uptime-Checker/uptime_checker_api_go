package watchdog

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/samber/lo"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
)

const (
	allGoodStatusCode = "200..299"
)

func (w *WatchDog) Assert(source int32, property *string, comparison int32, value string, resp HitResponse) bool {
	assertionSource := resource.AssertionSource(source)
	assertionComparison := resource.AssertionComparison(comparison)

	switch assertionSource {
	case resource.AssertionSourceStatusCode:
		return assertStatusCode(assertionComparison, value, resp.StatusCode)
	case resource.AssertionSourceResponseTime:
		return assertResponseTime(assertionComparison, value, resp.Duration)
	case resource.AssertionSourceTextBody:
		return assertTextBody(assertionComparison, value, resp.Body)
	case resource.AssertionSourceHeaders:
		return assertHeader(assertionComparison, *property, value, resp.Headers)
	}

	return false
}

func assertStatusCode(assertionComparison resource.AssertionComparison, value string, statusCode int) bool {
	if value == allGoodStatusCode {
		return checkSuccessStatusCode(statusCode)
	}

	code, err := strconv.Atoi(value)
	if err != nil {
		sentry.CaptureException(err)
		return false
	}

	switch assertionComparison {
	case resource.AssertionComparisonEqual:
		return code == statusCode
	case resource.AssertionComparisonNotEqual:
		return code != statusCode
	case resource.AssertionComparisonGreater:
		return code < statusCode
	case resource.AssertionComparisonLesser:
		return code > statusCode
	}
	return false
}

func checkSuccessStatusCode(code int) bool {
	codes := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent,
		http.StatusMultiStatus,
		http.StatusAlreadyReported,
		http.StatusIMUsed,
	}

	return lo.Contains(codes, code)
}

func assertResponseTime(assertionComparison resource.AssertionComparison, value string, duration int64) bool {
	responseTime, err := strconv.Atoi(value)
	if err != nil {
		sentry.CaptureException(err)
		return false
	}

	savedResponseTime := int64(responseTime * 1000)
	if assertionComparison == resource.AssertionComparisonGreater {
		return savedResponseTime < duration
	} else if assertionComparison == resource.AssertionComparisonLesser {
		return savedResponseTime > duration
	}
	return false
}

func assertTextBody(assertionComparison resource.AssertionComparison, value string, body *string) bool {
	if body == nil {
		return false
	}
	responseBody := *body

	switch assertionComparison {
	case resource.AssertionComparisonEqual:
		return value == responseBody
	case resource.AssertionComparisonNotEqual:
		return value != responseBody
	case resource.AssertionComparisonContain:
		return strings.Contains(responseBody, value)
	case resource.AssertionComparisonNotContain:
		return !strings.Contains(responseBody, value)
	case resource.AssertionComparisonEmpty:
		return value == ""
	case resource.AssertionComparisonNotEmpty:
		return value != ""
	}
	return false
}

func assertHeader(
	assertionComparison resource.AssertionComparison,
	property, value string,
	headers map[string]string,
) bool {
	headerValue, ok := headers[property]
	if !ok {
		return false
	}
	switch assertionComparison {
	case resource.AssertionComparisonEqual:
		return value == headerValue
	case resource.AssertionComparisonNotEqual:
		return value != headerValue
	case resource.AssertionComparisonContain:
		return strings.Contains(headerValue, value)
	case resource.AssertionComparisonNotContain:
		return !strings.Contains(headerValue, value)
	case resource.AssertionComparisonEmpty:
		return value == ""
	case resource.AssertionComparisonNotEmpty:
		return value != ""
	}
	return false
}
