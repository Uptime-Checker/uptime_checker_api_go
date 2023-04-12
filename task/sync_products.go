package task

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
)

type SyncProductsTask struct {
	productService *service.ProductService
}

func NewSyncProductsTask(productService *service.ProductService) *SyncProductsTask {
	return &SyncProductsTask{productService: productService}
}

func (s SyncProductsTask) Do(ctx context.Context, tx *sql.Tx) {
	tid := pkg.GetTracingID(ctx)

	lgr.Print(tid, 1, "running SyncProductsTask")
	billingProducts := infra.ListProductsWithPrices()
	for _, billingProduct := range billingProducts {
		product, err := s.productService.Add(
			ctx,
			tx,
			billingProduct.Name,
			billingProduct.Description,
			billingProduct.ID,
			billingProduct.Metadata["tier"],
		)
		if err != nil {
			sentry.CaptureException(errors.Newf("failed to add product %s, err: %w", billingProduct.Name, err))
			return
		}
	}
}
