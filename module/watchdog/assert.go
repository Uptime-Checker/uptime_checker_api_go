package watchdog

import (
	"net/http"
	"strconv"

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
	}

	return false
}

func assertStatusCode(assertionComparison resource.AssertionComparison, value string, statusCode int) bool {
	if value == allGoodStatusCode {
		return checkSuccessStatusCode(statusCode)
	}

	code, _ := strconv.Atoi(value)
	if assertionComparison == resource.AssertionComparisonEqual {
		if code == statusCode {
			return true
		}
	} else if assertionComparison == resource.AssertionComparisonNotEqual {
		if code == statusCode {
			return false
		}
	} else if assertionComparison == resource.AssertionComparisonGreater {
		if statusCode > code {
			return true
		}
	} else if assertionComparison == resource.AssertionComparisonLesser {
		if code < statusCode {
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
