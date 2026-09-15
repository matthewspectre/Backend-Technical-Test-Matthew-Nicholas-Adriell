package inventory

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	entity "be_evindo/internal/entity/inventory"
	usecase "be_evindo/internal/usecase/inventory"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
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

func inventoryResponse(data *entity.Inventory) InventoryResponse {
	return InventoryResponse{
		ID: data.ID, ProductID: data.ProductID, ProductName: data.ProductName,
		WarehouseID: data.WarehouseID, WarehouseName: data.WarehouseName,
		Stock: data.Stock, CreatedAt: data.CreatedAt.Format(responseTimeLayout),
		UpdatedAt: formatResponseTime(data.UpdatedAt),
	}
}

func isDuplicateInventoryError(err error) bool {
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return false
	}
	return postgresError.Code == "23505" && postgresError.ConstraintName == "inventory_product_warehouse_unique"
}

func (handler *Handler) Create(context *gin.Context) {
	var request InventoryCreateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data := &entity.Inventory{ProductID: request.ProductID, WarehouseID: request.WarehouseID, Stock: request.Stock}
	if err := handler.uc.Create(context.Request.Context(), data); err != nil {
		handler.writeError(context, err)
		return
	}
	createdInventory, err := handler.uc.FindByID(context.Request.Context(), data.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if createdInventory == nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "inventory could not be loaded after creation"})
		return
	}
	response := inventoryResponse(createdInventory)
	response.UpdatedAt = nil
	context.JSON(http.StatusCreated, response)
}

func (handler *Handler) GetAll(context *gin.Context) {
	data, err := handler.uc.FindAll(context.Request.Context())
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, handler.listResponse(data))
}

func (handler *Handler) GetByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid inventory id"})
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	if data == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
		return
	}
	context.JSON(http.StatusOK, inventoryResponse(data))
}

func (handler *Handler) GetByProduct(context *gin.Context) {
	productID, err := strconv.Atoi(context.Param("product_id"))
	if err != nil || productID <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	data, err := handler.uc.FindByProduct(context.Request.Context(), productID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, handler.listResponse(data))
}

func (handler *Handler) GetByWarehouse(context *gin.Context) {
	warehouseID, err := strconv.Atoi(context.Param("warehouse_id"))
	if err != nil || warehouseID <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse id"})
		return
	}
	data, err := handler.uc.FindByWarehouse(context.Request.Context(), warehouseID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, handler.listResponse(data))
}

func (handler *Handler) Update(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid inventory id"})
		return
	}
	var request InventoryUpdateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := make(map[string]interface{})
	if request.ProductID != nil {
		updates["product_id"] = *request.ProductID
	}
	if request.WarehouseID != nil {
		updates["warehouse_id"] = *request.WarehouseID
	}
	if request.Stock != nil {
		updates["stock"] = *request.Stock
	}
	if err := handler.uc.Update(context.Request.Context(), id, updates); err != nil {
		handler.writeError(context, err)
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, inventoryResponse(data))
}

func (handler *Handler) Delete(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid inventory id"})
		return
	}
	if err := handler.uc.Delete(context.Request.Context(), id); err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Inventory has been successfully deleted.", "id": id})
}

func (handler *Handler) listResponse(data []*entity.Inventory) InventoryListResponse {
	response := make([]InventoryResponse, 0, len(data))
	for _, item := range data {
		response = append(response, inventoryResponse(item))
	}
	return InventoryListResponse{Data: response, TotalData: len(response)}
}

func (handler *Handler) writeError(context *gin.Context, err error) {
	if errors.Is(err, usecase.ErrInactiveProduct) {
		context.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code":    "INVENTORY_PRODUCT_INACTIVE",
			"message": "Inactive product cannot be used for new inventory.",
		}})
		return
	}
	if errors.Is(err, usecase.ErrInactiveWarehouse) {
		context.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code":    "INVENTORY_WAREHOUSE_INACTIVE",
			"message": "Inactive warehouse cannot be used for new inventory.",
		}})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
		return
	}
	if isDuplicateInventoryError(err) {
		context.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code":    "INVENTORY_ALREADY_EXISTS",
			"message": "Inventory for this product and warehouse already exists.",
		}})
		return
	}
	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
