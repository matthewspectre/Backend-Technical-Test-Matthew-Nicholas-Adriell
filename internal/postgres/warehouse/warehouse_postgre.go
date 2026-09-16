package warehouse

import (
	"context"
	"time"

	entity "be_evindo/internal/entity/warehouse"
	model "be_evindo/internal/model/warehouse"
	repo "be_evindo/internal/repository/warehouse"

	"gorm.io/gorm"
)

type RepositoryPostgre struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repo.WarehouseRepository {
	return &RepositoryPostgre{db: db}
}

func toModel(entityWarehouse *entity.Warehouse) *model.WarehouseModel {
	if entityWarehouse == nil {
		return nil
	}
	return &model.WarehouseModel{
		ID:        entityWarehouse.ID,
		Code:      entityWarehouse.Code,
		Name:      entityWarehouse.Name,
		Location:  entityWarehouse.Location,
		IsActive:  entityWarehouse.IsActive,
		CreatedAt: entityWarehouse.CreatedAt,
		UpdatedAt: entityWarehouse.UpdatedAt,
	}
}

func toEntity(warehouseModel *model.WarehouseModel) *entity.Warehouse {
	if warehouseModel == nil {
		return nil
	}
	return &entity.Warehouse{
		ID:        warehouseModel.ID,
		Code:      warehouseModel.Code,
		Name:      warehouseModel.Name,
		Location:  warehouseModel.Location,
		IsActive:  warehouseModel.IsActive,
		CreatedAt: warehouseModel.CreatedAt,
		UpdatedAt: warehouseModel.UpdatedAt,
	}
}

func (repository *RepositoryPostgre) Create(ctx context.Context, entityWarehouse *entity.Warehouse) error {
	warehouseModel := toModel(entityWarehouse)
	if err := repository.db.WithContext(ctx).Create(warehouseModel).Error; err != nil {
		return err
	}
	if entityWarehouse != nil {
		entityWarehouse.ID = warehouseModel.ID
	}
	return nil
}

func (repository *RepositoryPostgre) FindByID(ctx context.Context, id int) (*entity.Warehouse, error) {
	warehouseModel := &model.WarehouseModel{}
	if err := repository.db.WithContext(ctx).First(warehouseModel, id).Error; err != nil {
		return nil, err
	}
	return toEntity(warehouseModel), nil
}

func (repository *RepositoryPostgre) FindAll(ctx context.Context) ([]*entity.Warehouse, error) {
	var warehouseModels []*model.WarehouseModel
	if err := repository.db.WithContext(ctx).
		Where("is_active = ?", 1).
		Order("id DESC").
		Find(&warehouseModels).Error; err != nil {
		return nil, err
	}

	warehouses := make([]*entity.Warehouse, 0, len(warehouseModels))
	for _, warehouseModel := range warehouseModels {
		warehouses = append(warehouses, toEntity(warehouseModel))
	}
	return warehouses, nil
}

func (repository *RepositoryPostgre) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	result := repository.db.WithContext(ctx).
		Model(&model.WarehouseModel{}).
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
