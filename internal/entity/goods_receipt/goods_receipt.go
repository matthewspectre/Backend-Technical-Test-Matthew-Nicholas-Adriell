package goods_receipt

import "time"

type GoodsReceipt struct {
	ID              int64               `json:"id"`
	ReceiptNumber   string              `json:"receipt_number"`
	PurchaseOrderID int64               `json:"purchase_order_id"`
	WarehouseID     int                 `json:"warehouse_id"`
	ReceivedBy      int64               `json:"received_by"`
	Status          string              `json:"status"`
	ReceivedAt      time.Time           `json:"received_at"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       *time.Time          `json:"updated_at"`
	Items           []*GoodsReceiptItem `json:"items"`
}

type GoodsReceiptItem struct {
	ID                  int64 `json:"id"`
	GoodsReceiptID      int64 `json:"goods_receipt_id"`
	PurchaseOrderItemID int64 `json:"purchase_order_item_id"`
	ProductID           int   `json:"product_id"`
	ReceivedQuantity    int   `json:"received_quantity"`
}
