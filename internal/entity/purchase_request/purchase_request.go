package purchase_request

import "time"

type PurchaseRequest struct {
	ID            int64                  `json:"id"`
	RequestNumber string                 `json:"request_number"`
	WarehouseID   int                    `json:"warehouse_id"`
	WarehouseName string                 `json:"warehouse_name"`
	RequestedBy   int64                  `json:"requested_by"`
	RequesterName string                 `json:"requester_name"`
	Status        string                 `json:"status"`
	Items         []*PurchaseRequestItem `json:"items"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at"`
}

type PurchaseRequestItem struct {
	ID                int64  `json:"id"`
	PurchaseRequestID int64  `json:"purchase_request_id"`
	ProductID         int    `json:"product_id"`
	ProductName       string `json:"product_name"`
	Unit              string `json:"unit"`
	Quantity          int    `json:"quantity"`
}
