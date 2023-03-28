package watchdog

import (
	"net/http"
	"strconv"

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
		if code == statusCode {
			return true
		}
	} else if assertionComparison == resource.AssertionComparisonNotEqual {
		if code == statusCode {
			return false
		}
	} else if assertionComparison == resource.AssertionComparisonGreater {
		if code < statusCode {
			return true
		}
	} else if assertionComparison == resource.AssertionComparisonLesser {
		if code > statusCode {
			return true
		}
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
		if savedResponseTime < duration {
			return true
		}
	} else if assertionComparison == resource.AssertionComparisonLesser {
		if savedResponseTime > duration {
			return true
		}
	}
	return false
}
