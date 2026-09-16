package inventory

import (
	"context"
	"errors"
	"time"

	entity "be_evindo/internal/entity/inventory"
	repo "be_evindo/internal/repository/inventory"
	productrepo "be_evindo/internal/repository/product"
	warehouserepo "be_evindo/internal/repository/warehouse"
)

type Usecase interface {
	Create(ctx context.Context, data *entity.Inventory) error
	FindAll(ctx context.Context) ([]*entity.Inventory, error)
	FindByID(ctx context.Context, id int) (*entity.Inventory, error)
	FindByProduct(ctx context.Context, productID int) ([]*entity.Inventory, error)
	FindByWarehouse(ctx context.Context, warehouseID int) ([]*entity.Inventory, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
}

type usecase struct {
	repository          repo.InventoryRepository
	warehouseRepository warehouserepo.WarehouseRepository
	productRepository   productrepo.ProductRepository
}

var ErrInactiveWarehouse = errors.New("inactive warehouse cannot be used for inventory")
var ErrInactiveProduct = errors.New("inactive product cannot be used for inventory")

func NewUsecase(repository repo.InventoryRepository, warehouseRepository warehouserepo.WarehouseRepository, productRepository productrepo.ProductRepository) Usecase {
	return &usecase{repository: repository, warehouseRepository: warehouseRepository, productRepository: productRepository}
}

func (usecase *usecase) Create(ctx context.Context, data *entity.Inventory) error {
	if data == nil {
		return errors.New("inventory is required")
	}
	if data.ProductID <= 0 || data.WarehouseID <= 0 {
		return errors.New("product_id and warehouse_id must be greater than zero")
	}
	if data.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	warehouse, err := usecase.warehouseRepository.FindByID(ctx, data.WarehouseID)
	if err != nil {
		return err
	}
	if warehouse == nil || warehouse.IsActive != 1 {
		return ErrInactiveWarehouse
	}
	product, err := usecase.productRepository.FindByID(ctx, data.ProductID)
	if err != nil {
		return err
	}
	if product == nil || product.IsActive != 1 {
		return ErrInactiveProduct
	}
	return usecase.repository.Create(ctx, data)
}

func (usecase *usecase) FindAll(ctx context.Context) ([]*entity.Inventory, error) {
	return usecase.repository.FindAll(ctx)
}

func (usecase *usecase) FindByID(ctx context.Context, id int) (*entity.Inventory, error) {
	if id <= 0 {
		return nil, errors.New("inventory id must be greater than zero")
	}
	return usecase.repository.FindByID(ctx, id)
}

func (usecase *usecase) FindByProduct(ctx context.Context, productID int) ([]*entity.Inventory, error) {
	if productID <= 0 {
		return nil, errors.New("product_id must be greater than zero")
	}
	return usecase.repository.FindByProduct(ctx, productID)
}

func (usecase *usecase) FindByWarehouse(ctx context.Context, warehouseID int) ([]*entity.Inventory, error) {
	if warehouseID <= 0 {
		return nil, errors.New("warehouse_id must be greater than zero")
	}
	return usecase.repository.FindByWarehouse(ctx, warehouseID)
}

func (usecase *usecase) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	if id <= 0 {
		return errors.New("inventory id must be greater than zero")
	}
	if len(updates) == 0 {
		return errors.New("at least one field is required")
	}
	for _, field := range []string{"product_id", "warehouse_id"} {
		if value, exists := updates[field]; exists {
			identifier, ok := value.(int)
			if !ok || identifier <= 0 {
				return errors.New(field + " must be greater than zero")
			}
		}
	}
	if value, exists := updates["stock"]; exists {
		stock, ok := value.(int)
		if !ok || stock < 0 {
			return errors.New("stock cannot be negative")
		}
	}
	if warehouseID, exists := updates["warehouse_id"]; exists {
		warehouse, err := usecase.warehouseRepository.FindByID(ctx, warehouseID.(int))
		if err != nil {
			return err
		}
		if warehouse == nil || warehouse.IsActive != 1 {
			return errors.New("inactive warehouses cannot be used for inventory transactions")
		}
	}
	updates["updated_at"] = time.Now()
	return usecase.repository.Update(ctx, id, updates)
}

func (usecase *usecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("inventory id must be greater than zero")
	}
	return usecase.repository.Delete(ctx, id)
}
