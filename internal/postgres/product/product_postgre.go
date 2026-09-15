package product

import (
	"context"
	"time"

	entity "be_evindo/internal/entity/product"
	model "be_evindo/internal/model/product"
	repo "be_evindo/internal/repository/product"

	"gorm.io/gorm"
)

type RepositoryPostgre struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repo.ProductRepository {
	return &RepositoryPostgre{db: db}
}

func toModel(entityProduct *entity.Product) *model.ProductModel {
	if entityProduct == nil {
		return nil
	}
	return &model.ProductModel{
		ID:        entityProduct.ID,
		SKU:       entityProduct.SKU,
		Name:      entityProduct.Name,
		Unit:      entityProduct.Unit,
		IsActive:  entityProduct.IsActive,
		CreatedAt: entityProduct.CreatedAt,
		UpdatedAt: entityProduct.UpdatedAt,
	}
}

func toEntity(productModel *model.ProductModel) *entity.Product {
	if productModel == nil {
		return nil
	}
	return &entity.Product{
		ID:        productModel.ID,
		SKU:       productModel.SKU,
		Name:      productModel.Name,
		Unit:      productModel.Unit,
		IsActive:  productModel.IsActive,
		CreatedAt: productModel.CreatedAt,
		UpdatedAt: productModel.UpdatedAt,
	}
}

func (repository *RepositoryPostgre) Create(ctx context.Context, entityProduct *entity.Product) error {
	productModel := toModel(entityProduct)
	if err := repository.db.WithContext(ctx).Create(productModel).Error; err != nil {
		return err
	}
	if entityProduct != nil {
		entityProduct.ID = productModel.ID
	}
	return nil
}

func (repository *RepositoryPostgre) FindByID(ctx context.Context, id int) (*entity.Product, error) {
	productModel := &model.ProductModel{}
	if err := repository.db.WithContext(ctx).First(productModel, id).Error; err != nil {
		return nil, err
	}
	return toEntity(productModel), nil
}

func (repository *RepositoryPostgre) FindAll(ctx context.Context) ([]*entity.Product, error) {
	var productModels []*model.ProductModel
	if err := repository.db.WithContext(ctx).
		Where("is_active = ?", 1).
		Order("id DESC").
		Find(&productModels).Error; err != nil {
		return nil, err
	}

	products := make([]*entity.Product, 0, len(productModels))
	for _, productModel := range productModels {
		products = append(products, toEntity(productModel))
	}
	return products, nil
}

func (repository *RepositoryPostgre) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	result := repository.db.WithContext(ctx).
		Model(&model.ProductModel{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repository *RepositoryPostgre) Delete(ctx context.Context, id int) error {
	return repository.Update(ctx, id, map[string]interface{}{
		"is_active":  0,
		"updated_at": time.Now(),
	})
}
