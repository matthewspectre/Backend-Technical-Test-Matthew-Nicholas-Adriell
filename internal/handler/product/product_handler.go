package product

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	entity "be_evindo/internal/entity/product"
	usecase "be_evindo/internal/usecase/product"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Handler struct {
	uc usecase.Usecase
}

const responseTimeLayout = "2006-01-02 15:04:05"

func formatResponseTime(value *time.Time) *string {
	if value == nil || value.IsZero() {
		return nil
	}
	formatted := value.Format(responseTimeLayout)
	return &formatted
}

func productResponse(data *entity.Product) ProductResponse {
	return ProductResponse{
		ID:        data.ID,
		SKU:       data.SKU,
		Name:      data.Name,
		Unit:      data.Unit,
		IsActive:  data.IsActive,
		CreatedAt: data.CreatedAt.Format(responseTimeLayout),
		UpdatedAt: formatResponseTime(data.UpdatedAt),
	}
}

func isDuplicateSKUError(err error) bool {
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return false
	}
	return postgresError.Code == "23505" &&
		(postgresError.ConstraintName == "product_unique" || postgresError.ColumnName == "sku")
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (handler *Handler) Create(context *gin.Context) {
	var request ProductCreateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := &entity.Product{
		SKU:      request.SKU,
		Name:     request.Name,
		Unit:     request.Unit,
		IsActive: request.IsActive,
	}
	if data.IsActive == 0 {
		data.IsActive = 1
	}

	if err := handler.uc.Create(context.Request.Context(), data); err != nil {
		if isDuplicateSKUError(err) {
			context.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "SKU_NOT_APPROVED",
					"message": "SKU already exists. Please provide a unique SKU.",
				},
			})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdProduct, err := handler.uc.FindByID(context.Request.Context(), data.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := productResponse(createdProduct)
	response.UpdatedAt = nil
	context.JSON(http.StatusCreated, response)
}

func (handler *Handler) GetAll(context *gin.Context) {
	data, err := handler.uc.FindAll(context.Request.Context())
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]ProductResponse, 0, len(data))
	for _, product := range data {
		response = append(response, productResponse(product))
	}
	context.JSON(http.StatusOK, ProductListResponse{
		Data:      response,
		TotalData: len(response),
	})
}

func (handler *Handler) Update(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var request ProductUpdateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if request.SKU != nil {
		updates["sku"] = *request.SKU
	}
	if request.Name != nil {
		updates["name"] = *request.Name
	}
	if request.Unit != nil {
		updates["unit"] = *request.Unit
	}
	if request.IsActive != nil {
		updates["is_active"] = *request.IsActive
	}

	if err := handler.uc.Update(context.Request.Context(), id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if isDuplicateSKUError(err) {
			context.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "SKU_NOT_APPROVED",
					"message": "SKU already exists. Please provide a unique SKU.",
				},
			})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedProduct, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if updatedProduct == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	context.JSON(http.StatusOK, productResponse(updatedProduct))
}

func (handler *Handler) Delete(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	if err := handler.uc.Delete(context.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Product has been successfully deactivated.",
		"id":      id,
	})
}

func (handler *Handler) GetByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	context.JSON(http.StatusOK, productResponse(data))
}
