package infra

import (
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/sourcegraph/conc/pool"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

func SetupBilling() {
	stripe.Key = config.App.StripeKey
	stripe.DefaultLeveledLogger = lgr.Zapper
}

func listProducts() []*stripe.Product {
	products := make([]*stripe.Product, 0)
	params := &stripe.ProductListParams{}
	i := product.List(params)
	for i.Next() {
		products = append(products, i.Product())
	}
	return products
}

func listPrices() []*stripe.Price {
	prices := make([]*stripe.Price, 0)
	params := &stripe.PriceListParams{}
	i := price.List(params)
	for i.Next() {
		prices = append(prices, i.Price())
	}
	return prices
}

func ListProductsWithPrices() []pkg.BillingProduct {
	var products []*stripe.Product
	var prices []*stripe.Price

	p := pool.New()
	p.Go(func() {
		products = listProducts()
	})
	p.Go(func() {
		prices = listPrices()
	})
	p.Wait()

	billingProducts := make([]pkg.BillingProduct, 0)

	for _, p := range products {
		if !p.Active {
			continue
		}
		billingProduct := pkg.BillingProduct{Product: p}
		billingProductPrices := make([]*stripe.Price, 0)
		for _, billingProductPrice := range prices {
			if billingProductPrice.Product.ID == p.ID {
				billingProductPrices = append(billingProductPrices, billingProductPrice)
			}
		}
		billingProduct.Prices = billingProductPrices
		billingProducts = append(billingProducts, billingProduct)
	}
	return billingProducts
}

func CreateCustomer(name *string, email string) (*stripe.Customer, error) {
	customerParams := &stripe.CustomerParams{
		Email: &email,
		Name:  name,
	}
	cus, err := customer.New(customerParams)
	if err != nil {
		sentry.CaptureException(errors.Newf("failed to create stripe customer for %s, err: %w", email, err))
		return nil, err
	}
	return cus, nil
}
