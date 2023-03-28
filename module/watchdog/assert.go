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
		if value == allGoodStatusCode {
			return checkSuccessStatusCode(resp.StatusCode)
		}

		code, _ := strconv.Atoi(value)
		if assertionComparison == resource.AssertionComparisonEqual {
			if code == resp.StatusCode {
				return true
			}
		} else if assertionComparison == resource.AssertionComparisonNotEqual {
			if code == resp.StatusCode {
				return false
			}
		} else if assertionComparison == resource.AssertionComparisonGreater {
			if resp.StatusCode > code {
				return true
			}
		} else if assertionComparison == resource.AssertionComparisonLesser {
			if code < resp.StatusCode {
				return true
			}
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

	return lo.Contains[int](codes, code)
}
