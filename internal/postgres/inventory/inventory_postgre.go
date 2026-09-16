package inventory

import (
	"context"
	"time"

	entity "be_evindo/internal/entity/inventory"
	model "be_evindo/internal/model/inventory"
	repo "be_evindo/internal/repository/inventory"

	"gorm.io/gorm"
)

type RepositoryPostgre struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repo.InventoryRepository {
	return &RepositoryPostgre{db: db}
}

func toModel(entityInventory *entity.Inventory) *model.InventoryModel {
	if entityInventory == nil {
		return nil
	}
	return &model.InventoryModel{
		ID:            entityInventory.ID,
		ProductID:     entityInventory.ProductID,
		ProductName:   entityInventory.ProductName,
		WarehouseID:   entityInventory.WarehouseID,
		WarehouseName: entityInventory.WarehouseName,
		Stock:         entityInventory.Stock,
		CreatedAt:     entityInventory.CreatedAt,
		UpdatedAt:     entityInventory.UpdatedAt,
	}
}

func toEntity(inventoryModel *model.InventoryModel) *entity.Inventory {
	if inventoryModel == nil {
		return nil
	}
	return &entity.Inventory{
		ID:            inventoryModel.ID,
		ProductID:     inventoryModel.ProductID,
		ProductName:   inventoryModel.ProductName,
		WarehouseID:   inventoryModel.WarehouseID,
		WarehouseName: inventoryModel.WarehouseName,
		Stock:         inventoryModel.Stock,
		CreatedAt:     inventoryModel.CreatedAt,
		UpdatedAt:     inventoryModel.UpdatedAt,
	}
}

func (repository *RepositoryPostgre) Create(ctx context.Context, entityInventory *entity.Inventory) error {
	inventoryModel := toModel(entityInventory)
	if err := repository.db.WithContext(ctx).Create(inventoryModel).Error; err != nil {
		return err
	}
	if entityInventory != nil {
		entityInventory.ID = inventoryModel.ID
	}
	return nil
}

func (repository *RepositoryPostgre) FindAll(ctx context.Context) ([]*entity.Inventory, error) {
	var inventoryModels []*model.InventoryModel
	if err := repository.db.WithContext(ctx).
		Table("inventory i").
		Select("i.*, p.name AS product_name, w.name AS warehouse_name").
		Joins("LEFT JOIN product p ON p.id = i.product_id").
		Joins("LEFT JOIN warehouse w ON w.id = i.warehouse_id").
		Order("id DESC").
		Find(&inventoryModels).Error; err != nil {
		return nil, err
	}

	inventories := make([]*entity.Inventory, 0, len(inventoryModels))
	for _, inventoryModel := range inventoryModels {
		inventories = append(inventories, toEntity(inventoryModel))
	}
	return inventories, nil
}

func (repository *RepositoryPostgre) FindByID(ctx context.Context, id int) (*entity.Inventory, error) {
	inventoryModel := &model.InventoryModel{}
	if err := repository.db.WithContext(ctx).
		Table("inventory i").
		Select("i.*, p.name AS product_name, w.name AS warehouse_name").
		Joins("LEFT JOIN product p ON p.id = i.product_id").
		Joins("LEFT JOIN warehouse w ON w.id = i.warehouse_id").
		Where("i.id = ?", id).
		First(inventoryModel).Error; err != nil {
		return nil, err
	}
	return toEntity(inventoryModel), nil
}

func (repository *RepositoryPostgre) FindByProduct(ctx context.Context, productID int) ([]*entity.Inventory, error) {
	var inventoryModels []*model.InventoryModel
	if err := repository.db.WithContext(ctx).
		Table("inventory i").
		Select("i.*, p.name AS product_name, w.name AS warehouse_name").
		Joins("LEFT JOIN product p ON p.id = i.product_id").
		Joins("LEFT JOIN warehouse w ON w.id = i.warehouse_id").
		Where("product_id = ?", productID).
		Order("warehouse_id ASC").
		Find(&inventoryModels).Error; err != nil {
		return nil, err
	}

	inventories := make([]*entity.Inventory, 0, len(inventoryModels))
	for _, inventoryModel := range inventoryModels {
		inventories = append(inventories, toEntity(inventoryModel))
	}
	return inventories, nil
}

func (repository *RepositoryPostgre) FindByWarehouse(ctx context.Context, warehouseID int) ([]*entity.Inventory, error) {
	var inventoryModels []*model.InventoryModel
	if err := repository.db.WithContext(ctx).
		Table("inventory i").
		Select("i.*, p.name AS product_name, w.name AS warehouse_name").
		Joins("LEFT JOIN product p ON p.id = i.product_id").
		Joins("LEFT JOIN warehouse w ON w.id = i.warehouse_id").
		Where("warehouse_id = ?", warehouseID).
		Order("product_id ASC").
		Find(&inventoryModels).Error; err != nil {
		return nil, err
	}

	inventories := make([]*entity.Inventory, 0, len(inventoryModels))
	for _, inventoryModel := range inventoryModels {
		inventories = append(inventories, toEntity(inventoryModel))
	}
	return inventories, nil
}

func (repository *RepositoryPostgre) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	result := repository.db.WithContext(ctx).
		Model(&model.InventoryModel{}).
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
	result := repository.db.WithContext(ctx).
		Delete(&model.InventoryModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
