package purchase_order

import "time"

type PurchaseOrderModel struct {
	ID                int64                     `gorm:"primaryKey;autoIncrement;column:id"`
	PONumber          string                    `gorm:"column:po_number;default:(-)"`
	PurchaseRequestID int64                     `gorm:"column:purchase_request_id"`
	SupplierID        int                       `gorm:"column:supplier_id"`
	SupplierName      string                    `gorm:"column:supplier_name;->"`
	WarehouseID       int                       `gorm:"column:warehouse_id"`
	WarehouseName     string                    `gorm:"column:warehouse_name;->"`
	Status            string                    `gorm:"column:status"`
	Items             []*PurchaseOrderItemModel `gorm:"-"`
	CreatedAt         time.Time                 `gorm:"column:created_at"`
	UpdatedAt         *time.Time                `gorm:"column:updated_at"`
}

type PurchaseOrderItemModel struct {
	ID               int64  `gorm:"primaryKey;autoIncrement;column:id"`
	PurchaseOrderID  int64  `gorm:"column:purchase_order_id"`
	ProductID        int    `gorm:"column:product_id"`
	ProductName      string `gorm:"column:product_name;->"`
	Unit             string `gorm:"column:unit;->"`
	OrderedQuantity  int    `gorm:"column:ordered_quantity"`
	ReceivedQuantity int    `gorm:"column:received_quantity"`
}

func (PurchaseOrderItemModel) TableName() string {
	return "purchase_order_items"
}

func (PurchaseOrderModel) TableName() string {
	return "purchase_orders"
}
