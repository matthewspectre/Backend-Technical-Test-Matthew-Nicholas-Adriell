package inventory

import "time"

type Inventory struct {
	ID            int        `json:"id"`
	ProductID     int        `json:"product_id"`
	ProductName   string     `json:"product_name"`
	WarehouseID   int        `json:"warehouse_id"`
	WarehouseName string     `json:"warehouse_name"`
	Stock         int        `json:"stock"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
