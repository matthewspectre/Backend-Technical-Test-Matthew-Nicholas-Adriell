package supplier

import (
	"context"
	"time"

	entity "be_evindo/internal/entity/supplier"
	model "be_evindo/internal/model/supplier"
	repo "be_evindo/internal/repository/supplier"

	"gorm.io/gorm"
)

type RepositoryPostgre struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repo.SupplierRepository {
	return &RepositoryPostgre{db: db}
}

func toModel(entitySupplier *entity.Supplier) *model.SupplierModel {
	if entitySupplier == nil {
		return nil
	}
	return &model.SupplierModel{
		ID:          entitySupplier.ID,
		CompanyName: entitySupplier.CompanyName,
		Name:        entitySupplier.Name,
		Email:       entitySupplier.Email,
		Phone:       entitySupplier.Phone,
		Address:     entitySupplier.Address,
		IsActive:    entitySupplier.IsActive,
		CreatedAt:   entitySupplier.CreatedAt,
		UpdatedAt:   entitySupplier.UpdatedAt,
	}
}

func toEntity(supplierModel *model.SupplierModel) *entity.Supplier {
	if supplierModel == nil {
		return nil
	}
	return &entity.Supplier{
		ID:          supplierModel.ID,
		CompanyName: supplierModel.CompanyName,
		Name:        supplierModel.Name,
		Email:       supplierModel.Email,
		Phone:       supplierModel.Phone,
		Address:     supplierModel.Address,
		IsActive:    supplierModel.IsActive,
		CreatedAt:   supplierModel.CreatedAt,
		UpdatedAt:   supplierModel.UpdatedAt,
	}
}

func (repository *RepositoryPostgre) Create(ctx context.Context, entitySupplier *entity.Supplier) error {
	supplierModel := toModel(entitySupplier)
	if err := repository.db.WithContext(ctx).Create(supplierModel).Error; err != nil {
		return err
	}
	if entitySupplier != nil {
		entitySupplier.ID = supplierModel.ID
	}
	return nil
}

func (repository *RepositoryPostgre) FindByID(ctx context.Context, id int) (*entity.Supplier, error) {
	supplierModel := &model.SupplierModel{}
	if err := repository.db.WithContext(ctx).First(supplierModel, id).Error; err != nil {
		return nil, err
	}
	return toEntity(supplierModel), nil
}

func (repository *RepositoryPostgre) FindAll(ctx context.Context) ([]*entity.Supplier, error) {
	var supplierModels []*model.SupplierModel
	if err := repository.db.WithContext(ctx).
		Where("is_active = ?", 1).
		Order("id DESC").
		Find(&supplierModels).Error; err != nil {
		return nil, err
	}

	suppliers := make([]*entity.Supplier, 0, len(supplierModels))
	for _, supplierModel := range supplierModels {
		suppliers = append(suppliers, toEntity(supplierModel))
	}
	return suppliers, nil
}

func (repository *RepositoryPostgre) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	result := repository.db.WithContext(ctx).
		Model(&model.SupplierModel{}).
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
