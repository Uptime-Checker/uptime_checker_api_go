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

	if assertionSource == resource.AssertionSourceStatusCode {
		return assertStatusCode(assertionComparison, value, resp.StatusCode)
	} else if assertionSource == resource.AssertionSourceResponseTime {
		return assertResponseTime(assertionComparison, value, resp.Duration)
	} else if assertionSource == resource.AssertionSourceTextBody {
		return assertTextBody(assertionComparison, value, resp.Body)
	} else if assertionSource == resource.AssertionSourceHeaders {
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

	if assertionComparison == resource.AssertionComparisonEqual {
		return code == statusCode
	} else if assertionComparison == resource.AssertionComparisonNotEqual {
		return code != statusCode
	} else if assertionComparison == resource.AssertionComparisonGreater {
		return code < statusCode
	} else if assertionComparison == resource.AssertionComparisonLesser {
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

	if assertionComparison == resource.AssertionComparisonEqual {
		return value == responseBody
	} else if assertionComparison == resource.AssertionComparisonNotEqual {
		return value != responseBody
	} else if assertionComparison == resource.AssertionComparisonContain {
		return strings.Contains(responseBody, value)
	} else if assertionComparison == resource.AssertionComparisonNotContain {
		return !strings.Contains(responseBody, value)
	} else if assertionComparison == resource.AssertionComparisonEmpty {
		return value == ""
	} else if assertionComparison == resource.AssertionComparisonNotEmpty {
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
	if assertionComparison == resource.AssertionComparisonEqual {
		return value == headerValue
	} else if assertionComparison == resource.AssertionComparisonNotEqual {
		return value != headerValue
	} else if assertionComparison == resource.AssertionComparisonContain {
		return strings.Contains(headerValue, value)
	} else if assertionComparison == resource.AssertionComparisonNotContain {
		return !strings.Contains(headerValue, value)
	} else if assertionComparison == resource.AssertionComparisonEmpty {
		return value == ""
	} else if assertionComparison == resource.AssertionComparisonNotEmpty {
		return value != ""
	}
	return false
}
