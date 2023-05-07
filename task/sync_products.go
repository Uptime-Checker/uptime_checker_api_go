package task

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/stripe/stripe-go/v74"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
)

type SyncProductsTask struct {
	planDomain     *domain.PlanDomain
	productService *service.ProductService
}

func NewSyncProductsTask(planDomain *domain.PlanDomain, productService *service.ProductService) *SyncProductsTask {
	return &SyncProductsTask{planDomain: planDomain, productService: productService}
}

func (s SyncProductsTask) Do(ctx context.Context) {
	tid := pkg.GetTracingID(ctx)

	if err := infra.Transaction(ctx, func(tx *sql.Tx) error {
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
				return errors.Newf("failed to add product %s, err: %w", billingProduct.Name, err)
			}

			for _, price := range billingProduct.Prices {
				plan := &model.Plan{
					ExternalID: &price.ID,
					ProductID:  &product.ID,
				}
				_, err := s.planDomain.Create(
					ctx,
					tx,
					plan,
					pkg.CentsToDollars(int(price.UnitAmount)),
					s.getPlantType(price.Recurring.Interval),
				)
				if err != nil {
					return errors.Newf("failed to add plan %s, err: %w", price.ID, err)
				}
			}
		}
		return nil
	}); err != nil {
		sentry.CaptureException(err)
	}
}

func (s SyncProductsTask) getPlantType(planType stripe.PriceRecurringInterval) resource.PlanType {
	switch planType {
	case stripe.PriceRecurringIntervalMonth:
		return resource.PlanTypeMonthly
	case stripe.PriceRecurringIntervalYear:
		return resource.PlanTypeYearly
	}
	return resource.PlanTypeMonthly
}
