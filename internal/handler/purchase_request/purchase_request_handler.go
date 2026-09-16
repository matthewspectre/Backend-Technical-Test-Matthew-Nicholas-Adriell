package purchase_request

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	entity "be_evindo/internal/entity/purchase_request"
	usecase "be_evindo/internal/usecase/purchase_request"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	uc usecase.Usecase
}

const responseTimeLayout = "2006-01-02 15:04:05"

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func formatResponseTime(value *time.Time) *string {
	if value == nil || value.IsZero() {
		return nil
	}
	formatted := value.Format(responseTimeLayout)
	return &formatted
}

func purchaseRequestResponse(data *entity.PurchaseRequest) PurchaseRequestResponse {
	items := make([]PurchaseRequestItemResponse, 0, len(data.Items))
	for _, item := range data.Items {
		items = append(items, PurchaseRequestItemResponse{
			ID: item.ID, ProductID: item.ProductID, ProductName: item.ProductName,
			Unit: item.Unit, Quantity: item.Quantity,
		})
	}
	return PurchaseRequestResponse{
		ID: data.ID, RequestNumber: data.RequestNumber, WarehouseID: data.WarehouseID,
		WarehouseName: data.WarehouseName, RequestedBy: data.RequestedBy,
		RequesterName: data.RequesterName, Status: data.Status,
		Items:     items,
		CreatedAt: data.CreatedAt.Format(responseTimeLayout),
		UpdatedAt: formatResponseTime(data.UpdatedAt),
	}
}

func (handler *Handler) Create(context *gin.Context) {
	role, exists := context.Get("user_role")
	if !exists || role != "USER" {
		context.JSON(http.StatusForbidden, gin.H{"error": "only users can create purchase requests"})
		return
	}
	userIDValue, exists := context.Get("user_id")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "user identity is missing"})
		return
	}
	userID, ok := userIDValue.(int64)
	if !ok || userID <= 0 {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user identity"})
		return
	}
	var request PurchaseRequestCreateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data := &entity.PurchaseRequest{
		RequestNumber: request.RequestNumber, WarehouseID: request.WarehouseID,
		RequestedBy: userID, Items: make([]*entity.PurchaseRequestItem, 0, len(request.Items)),
	}
	for _, item := range request.Items {
		data.Items = append(data.Items, &entity.PurchaseRequestItem{ProductID: item.ProductID, Quantity: item.Quantity})
	}
	if err := handler.uc.Create(context.Request.Context(), data); err != nil {
		handler.writeError(context, err)
		return
	}
	created, err := handler.uc.FindByID(context.Request.Context(), data.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, purchaseRequestResponse(created))
}

func (handler *Handler) GetAll(context *gin.Context) {
	data, err := handler.uc.FindAll(context.Request.Context(), context.Query("status"))
	if err != nil {
		handler.writeError(context, err)
		return
	}
	response := make([]PurchaseRequestResponse, 0, len(data))
	for _, item := range data {
		response = append(response, purchaseRequestResponse(item))
	}
	context.JSON(http.StatusOK, PurchaseRequestListResponse{Data: response, TotalData: len(response)})
}

func (handler *Handler) GetByID(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase request id"})
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, purchaseRequestResponse(data))
}

func (handler *Handler) Update(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase request id"})
		return
	}
	var request PurchaseRequestUpdateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var items []*entity.PurchaseRequestItem
	if request.Items != nil {
		items = make([]*entity.PurchaseRequestItem, 0, len(*request.Items))
		for _, item := range *request.Items {
			items = append(items, &entity.PurchaseRequestItem{ProductID: item.ProductID, Quantity: item.Quantity})
		}
	}
	if err := handler.uc.Update(context.Request.Context(), id, request.WarehouseID, request.Status, items); err != nil {
		handler.writeError(context, err)
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, purchaseRequestResponse(data))
}

func (handler *Handler) Approve(context *gin.Context) {
	role, exists := context.Get("user_role")
	if !exists || role != "APPROVER" {
		context.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_APPROVER_ONLY",
				"message": "Only users with APPROVER role can approve purchase requests.",
			},
		})
		return
	}
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase request id"})
		return
	}
	if err := handler.uc.Approve(context.Request.Context(), id); err != nil {
		handler.writeError(context, err)
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, purchaseRequestResponse(data))
}

func (handler *Handler) Reject(context *gin.Context) {
	role, exists := context.Get("user_role")
	if !exists || role != "APPROVER" {
		context.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_APPROVER_ONLY",
				"message": "Only users with APPROVER role can reject purchase requests.",
			},
		})
		return
	}
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase request id"})
		return
	}
	if err := handler.uc.Reject(context.Request.Context(), id); err != nil {
		handler.writeError(context, err)
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, purchaseRequestResponse(data))
}

func (handler *Handler) writeError(context *gin.Context, err error) {
	if errors.Is(err, usecase.ErrInactiveWarehouse) {
		context.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code":    "PURCHASE_REQUEST_WAREHOUSE_INACTIVE",
			"message": "Inactive warehouse cannot be used for new Purchase Request.",
		}})
		return
	}
	if errors.Is(err, usecase.ErrInactiveProduct) {
		context.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code":    "PURCHASE_REQUEST_PRODUCT_INACTIVE",
			"message": "Inactive product cannot be used for new Purchase Request.",
		}})
		return
	}
	if errors.Is(err, usecase.ErrApprovalRequiresSubmitted) {
		context.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_NOT_SUBMITTED",
				"message": "Only SUBMITTED purchase requests can be approved.",
			},
		})
		return
	}
	if errors.Is(err, usecase.ErrRejectionRequiresSubmitted) {
		context.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_NOT_SUBMITTED",
				"message": "Only SUBMITTED purchase requests can be rejected.",
			},
		})
		return
	}
	if errors.Is(err, usecase.ErrNotDraftNotEditable) {
		context.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_NOT_DRAFT",
				"message": "Only DRAFT purchase requests can be edited.",
			},
		})
		return
	}
	if errors.Is(err, usecase.ErrUpdateFieldsRequired) {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_UPDATE_FIELDS_REQUIRED",
				"message": "At least one field is required for update.",
			},
		})
		return
	}
	if errors.Is(err, usecase.ErrInvalidStatus) {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_STATUS_INVALID",
				"message": "Status must be DRAFT, SUBMITTED, APPROVED, or REJECTED.",
			},
		})
		return
	}
	if errors.Is(err, usecase.ErrStatusDataNotFound) {
		context.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_STATUS_DATA_NOT_FOUND",
				"message": "No purchase requests found for the selected status.",
			},
		})
		return
	}
	if errors.Is(err, usecase.ErrItemsRequired) {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_ITEMS_REQUIRED",
				"message": "At least one item is required.",
			},
		})
		return
	}
	if errors.Is(err, usecase.ErrInvalidItemQuantity) {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "PURCHASE_REQUEST_ITEM_QUANTITY_INVALID",
				"message": "Item quantity must be greater than zero.",
			},
		})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "purchase request not found"})
		return
	}
	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
