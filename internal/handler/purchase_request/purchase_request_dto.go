package purchase_request

type PurchaseRequestCreateRequest struct {
	RequestNumber string                             `json:"request_number"`
	WarehouseID   int                                `json:"warehouse_id"`
	Items         []PurchaseRequestItemCreateRequest `json:"items"`
}

type PurchaseRequestItemCreateRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type PurchaseRequestUpdateRequest struct {
	WarehouseID *int                                `json:"warehouse_id,omitempty"`
	Status      *string                             `json:"status,omitempty"`
	Items       *[]PurchaseRequestItemCreateRequest `json:"items,omitempty"`
}

type PurchaseRequestResponse struct {
	ID            int64                         `json:"id"`
	RequestNumber string                        `json:"request_number"`
	WarehouseID   int                           `json:"warehouse_id"`
	WarehouseName string                        `json:"warehouse_name"`
	RequestedBy   int64                         `json:"requested_by"`
	RequesterName string                        `json:"requester_name"`
	Status        string                        `json:"status"`
	Items         []PurchaseRequestItemResponse `json:"items"`
	CreatedAt     string                        `json:"created_at"`
	UpdatedAt     *string                       `json:"updated_at"`
}

type PurchaseRequestItemResponse struct {
	ID          int64  `json:"id"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Unit        string `json:"unit"`
	Quantity    int    `json:"quantity"`
}

type PurchaseRequestListResponse struct {
	Data      []PurchaseRequestResponse `json:"data"`
	TotalData int                       `json:"total_data"`
}
