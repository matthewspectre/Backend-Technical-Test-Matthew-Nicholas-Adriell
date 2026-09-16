package inventory

import "time"

type InventoryModel struct {
	ID            int        `gorm:"primaryKey;autoIncrement;column:id"`
	ProductID     int        `gorm:"column:product_id"`
	ProductName   string     `gorm:"column:product_name;->"`
	WarehouseID   int        `gorm:"column:warehouse_id"`
	WarehouseName string     `gorm:"column:warehouse_name;->"`
	Stock         int        `gorm:"column:stock"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
}

func (InventoryModel) TableName() string {
	return "inventory"
}
