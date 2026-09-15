package warehouse

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	entity "be_evindo/internal/entity/warehouse"
	usecase "be_evindo/internal/usecase/warehouse"

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

func warehouseResponse(data *entity.Warehouse) WarehouseResponse {
	return WarehouseResponse{
		ID: data.ID, Code: data.Code, Name: data.Name, Location: data.Location,
		IsActive: data.IsActive, CreatedAt: data.CreatedAt.Format(responseTimeLayout),
		UpdatedAt: formatResponseTime(data.UpdatedAt),
	}
}

func isDuplicateWarehouseCodeError(err error) bool {
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return false
	}
	return postgresError.Code == "23505" &&
		(postgresError.ConstraintName == "warehouse_code_unique" || postgresError.ColumnName == "code")
}

func (handler *Handler) Create(context *gin.Context) {
	var request WarehouseCreateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data := &entity.Warehouse{Code: request.Code, Name: request.Name, Location: request.Location, IsActive: request.IsActive}
	if data.IsActive == 0 {
		data.IsActive = 1
	}
	if err := handler.uc.Create(context.Request.Context(), data); err != nil {
		handler.writeError(context, err)
		return
	}
	createdWarehouse, err := handler.uc.FindByID(context.Request.Context(), data.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if createdWarehouse == nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "warehouse could not be loaded after creation"})
		return
	}
	response := warehouseResponse(createdWarehouse)
	response.UpdatedAt = nil
	context.JSON(http.StatusCreated, response)
}

func (handler *Handler) GetAll(context *gin.Context) {
	data, err := handler.uc.FindAll(context.Request.Context())
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := make([]WarehouseResponse, 0, len(data))
	for _, warehouse := range data {
		response = append(response, warehouseResponse(warehouse))
	}
	context.JSON(http.StatusOK, WarehouseListResponse{Data: response, TotalData: len(response)})
}

func (handler *Handler) GetByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse id"})
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
		return
	}
	context.JSON(http.StatusOK, warehouseResponse(data))
}

func (handler *Handler) Update(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse id"})
		return
	}
	var request WarehouseUpdateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if request.Code != nil {
		updates["code"] = *request.Code
	}
	if request.Name != nil {
		updates["name"] = *request.Name
	}
	if request.Location != nil {
		updates["location"] = *request.Location
	}
	if request.IsActive != nil {
		updates["is_active"] = *request.IsActive
	}
	if err := handler.uc.Update(context.Request.Context(), id, updates); err != nil {
		handler.writeError(context, err)
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
		return
	}
	context.JSON(http.StatusOK, warehouseResponse(data))
}

func (handler *Handler) Delete(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse id"})
		return
	}
	if err := handler.uc.Delete(context.Request.Context(), id); err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Warehouse has been successfully deactivated.", "id": id})
}

func (handler *Handler) writeError(context *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
		return
	}
	if isDuplicateWarehouseCodeError(err) {
		context.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code":    "WAREHOUSE_CODE_ALREADY_EXISTS",
			"message": "A warehouse with this code already exists.",
		}})
		return
	}
	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
