package infra

import (
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
)

func SetupBilling() {
	stripe.Key = config.App.StripeKey
	stripe.DefaultLeveledLogger = lgr.Zapper
}

func ListProducts() []*stripe.Product {
	products := make([]*stripe.Product, 0)
	params := &stripe.ProductListParams{}
	i := product.List(params)
	for i.Next() {
		products = append(products, i.Product())
	}
	return products
}

func ListPrices() []*stripe.Price {
	prices := make([]*stripe.Price, 0)
	params := &stripe.PriceListParams{}
	i := price.List(params)
	for i.Next() {
		prices = append(prices, i.Price())
	}
	return prices
}

func ListProductsWithPrices() {
	
}
