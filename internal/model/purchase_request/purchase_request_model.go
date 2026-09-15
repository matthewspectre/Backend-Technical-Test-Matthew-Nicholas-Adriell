package purchase_request

import "time"

type PurchaseRequestModel struct {
	ID            int64                       `gorm:"primaryKey;autoIncrement;column:id"`
	RequestNumber string                      `gorm:"column:request_number;default:(-)"`
	WarehouseID   int                         `gorm:"column:warehouse_id"`
	WarehouseName string                      `gorm:"column:warehouse_name;->"`
	RequestedBy   int64                       `gorm:"column:requested_by"`
	RequesterName string                      `gorm:"column:requester_name;->"`
	Status        string                      `gorm:"column:status"`
	Items         []*PurchaseRequestItemModel `gorm:"-"`
	CreatedAt     time.Time                   `gorm:"column:created_at"`
	UpdatedAt     *time.Time                  `gorm:"column:updated_at"`
}

type PurchaseRequestItemModel struct {
	ID                int64  `gorm:"primaryKey;autoIncrement;column:id"`
	PurchaseRequestID int64  `gorm:"column:purchase_request_id"`
	ProductID         int    `gorm:"column:product_id"`
	ProductName       string `gorm:"column:product_name;->"`
	Unit              string `gorm:"column:unit;->"`
	Quantity          int    `gorm:"column:quantity"`
}

func (PurchaseRequestItemModel) TableName() string {
	return "purchase_request_items"
}

func (PurchaseRequestModel) TableName() string {
	return "purchase_requests"
}
