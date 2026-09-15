package goods_receipt

import "time"

type GoodsReceiptModel struct {
	ID              int64                    `gorm:"primaryKey;autoIncrement;column:id"`
	ReceiptNumber   string                   `gorm:"column:receipt_number"`
	PurchaseOrderID int64                    `gorm:"column:purchase_order_id"`
	WarehouseID     int                      `gorm:"column:warehouse_id"`
	ReceivedBy      int64                    `gorm:"column:received_by"`
	Status          string                   `gorm:"column:status"`
	ReceivedAt      time.Time                `gorm:"column:received_at"`
	CreatedAt       time.Time                `gorm:"column:created_at"`
	UpdatedAt       *time.Time               `gorm:"column:updated_at"`
	Items           []*GoodsReceiptItemModel `gorm:"-"`
}

type GoodsReceiptItemModel struct {
	ID                  int64 `gorm:"primaryKey;autoIncrement;column:id"`
	GoodsReceiptID      int64 `gorm:"column:goods_receipt_id"`
	PurchaseOrderItemID int64 `gorm:"column:purchase_order_item_id"`
	ProductID           int   `gorm:"column:product_id"`
	ReceivedQuantity    int   `gorm:"column:received_quantity"`
}

func (GoodsReceiptModel) TableName() string     { return "goods_receipts" }
func (GoodsReceiptItemModel) TableName() string { return "goods_receipt_items" }
