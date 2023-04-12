package service

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type ProductService struct {
	productDomain *domain.ProductDomain
}

func NewProductService(productDomain *domain.ProductDomain) *ProductService {
	return &ProductService{productDomain: productDomain}
}

func (p *ProductService) Add(
	ctx context.Context,
	tx *sql.Tx, name, description, externalID, tier string,
) (*model.Product, error) {
	product, err := p.productDomain.Get(ctx, name)
	if err != nil {
		// Create
		product = &model.Product{
			Name:        name,
			Description: &description,
			ExternalID:  &externalID,
		}
		return p.productDomain.Create(ctx, tx, product, resource.GetProductTier(tier))
	}

	return product, nil
}
