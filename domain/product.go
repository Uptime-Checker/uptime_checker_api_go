package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type ProductDomain struct{}

func NewProductDomain() *ProductDomain {
	return &ProductDomain{}
}

func (p *ProductDomain) ListProducts(ctx context.Context) ([]model.Product, error) {
	stmt := SELECT(Product.AllColumns).FROM(Product)

	var products []model.Product
	err := stmt.QueryContext(ctx, infra.DB, &products)
	return products, err
}

func (p *ProductDomain) Get(ctx context.Context, name string) (*model.Product, error) {
	stmt := SELECT(Product.AllColumns).FROM(Product).WHERE(Product.Name.EQ(String(name))).LIMIT(1)

	product := &model.Product{}
	err := stmt.QueryContext(ctx, infra.DB, product)
	return product, err
}

func (p *ProductDomain) Create(
	ctx context.Context, tx *sql.Tx, product *model.Product,
	productTier resource.ProductTier,
) (*model.Product, error) {
	if !productTier.Valid() {
		return nil, constant.ErrInvalidProductTier
	}
	product.Tier = int32(productTier)
	insertStmt := Product.INSERT(Product.MutableColumns.Except(Product.InsertedAt, Product.UpdatedAt)).
		MODEL(product).RETURNING(Product.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, product)
	return product, err
}
