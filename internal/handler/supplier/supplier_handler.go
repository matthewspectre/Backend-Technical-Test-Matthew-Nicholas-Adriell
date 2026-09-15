package supplier

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	entity "be_evindo/internal/entity/supplier"
	usecase "be_evindo/internal/usecase/supplier"

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

func supplierResponse(data *entity.Supplier) SupplierResponse {
	return SupplierResponse{
		ID:          data.ID,
		CompanyName: data.CompanyName,
		Name:        data.Name,
		Email:       data.Email,
		Phone:       data.Phone,
		Address:     data.Address,
		IsActive:    data.IsActive,
		CreatedAt:   data.CreatedAt.Format(responseTimeLayout),
		UpdatedAt:   formatResponseTime(data.UpdatedAt),
	}
}

func isDuplicateSupplierEmailError(err error) bool {
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return false
	}
	return postgresError.Code == "23505" &&
		(postgresError.ConstraintName == "supplier_email_unique" || postgresError.ColumnName == "email")
}

func (handler *Handler) Create(context *gin.Context) {
	var request SupplierCreateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data := &entity.Supplier{
		CompanyName: request.CompanyName,
		Name:        request.Name,
		Email:       request.Email,
		Phone:       request.Phone,
		Address:     request.Address,
		IsActive:    request.IsActive,
	}
	if data.IsActive == 0 {
		data.IsActive = 1
	}
	if err := handler.uc.Create(context.Request.Context(), data); err != nil {
		handler.writeError(context, err)
		return
	}
	createdSupplier, err := handler.uc.FindByID(context.Request.Context(), data.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if createdSupplier == nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "supplier could not be loaded after creation"})
		return
	}
	response := supplierResponse(createdSupplier)
	response.UpdatedAt = nil
	context.JSON(http.StatusCreated, response)
}

func (handler *Handler) GetAll(context *gin.Context) {
	data, err := handler.uc.FindAll(context.Request.Context())
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := make([]SupplierResponse, 0, len(data))
	for _, supplier := range data {
		response = append(response, supplierResponse(supplier))
	}
	context.JSON(http.StatusOK, SupplierListResponse{Data: response, TotalData: len(response)})
}

func (handler *Handler) GetByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		return
	}
	data, err := handler.uc.FindByID(context.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
		return
	}
	context.JSON(http.StatusOK, supplierResponse(data))
}

func (handler *Handler) Update(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		return
	}
	var request SupplierUpdateRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if request.CompanyName != nil {
		updates["company_name"] = *request.CompanyName
	}
	if request.Name != nil {
		updates["name"] = *request.Name
	}
	if request.Email != nil {
		updates["email"] = *request.Email
	}
	if request.Phone != nil {
		updates["phone"] = *request.Phone
	}
	if request.Address != nil {
		updates["address"] = *request.Address
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
		context.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
		return
	}
	context.JSON(http.StatusOK, supplierResponse(data))
}

func (handler *Handler) Delete(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid supplier id"})
		return
	}
	if err := handler.uc.Delete(context.Request.Context(), id); err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Supplier has been successfully deactivated.", "id": id})
}

func (handler *Handler) writeError(context *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
		return
	}
	if isDuplicateSupplierEmailError(err) {
		context.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code":    "SUPPLIER_EMAIL_ALREADY_EXISTS",
			"message": "A supplier with this email address already exists.",
		}})
		return
	}
	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
