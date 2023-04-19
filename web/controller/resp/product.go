package resp

import (
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

type Product struct {
	Popular bool
	pkg.ProductWithPlansAndFeatures
}
