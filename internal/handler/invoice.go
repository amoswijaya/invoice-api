package handler

import (
	"fmt"
	"invoice-api/internal/db"
	"invoice-api/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type invoiceReq struct {
	InvoiceNumber string              `json:"invoiceNumber" binding:"required,min=1"`
	Date          string              `json:"date" binding:"required"`
	FromName      string              `json:"fromName" binding:"required"`
	FromEmail     string              `json:"fromEmail" binding:"required,email"`
	ToName        string              `json:"toName" binding:"required"`
	ToEmail       string              `json:"toEmail" binding:"required,email"`
	TaxRate       float64             `json:"taxRate" binding:"min=0,max=100"`
	Subtotal      float64             `json:"subtotal" binding:"required,min=0"`
	TaxAmount     float64             `json:"taxAmount" binding:"required,min=0"`
	Total         float64             `json:"total" binding:"required,min=0"`
	Items         []model.InvoiceItem `json:"items" binding:"required,min=1,dive"`
}

// Helper function untuk handle error responses
func handleError(c *gin.Context, statusCode int, message string, err error) {
	response := gin.H{"error": message}
	if err != nil && gin.Mode() == gin.DebugMode {
		response["details"] = err.Error()
	}
	c.JSON(statusCode, response)
}

// Helper function untuk validate date
func validateDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	return date, nil
}

// Helper function untuk validate invoice calculations
func validateInvoiceCalculations(req *invoiceReq) error {
	// Calculate expected totals
	expectedSubtotal := 0.0
	for _, item := range req.Items {
		expectedSubtotal += item.Amount
	}
	
	expectedTaxAmount := expectedSubtotal * (req.TaxRate / 100)
	expectedTotal := expectedSubtotal + expectedTaxAmount
	
	// Allow small floating point differences (0.01)
	tolerance := 0.01
	
	if abs(req.Subtotal-expectedSubtotal) > tolerance {
		return fmt.Errorf("subtotal mismatch: expected %.2f, got %.2f", expectedSubtotal, req.Subtotal)
	}
	
	if abs(req.TaxAmount-expectedTaxAmount) > tolerance {
		return fmt.Errorf("tax amount mismatch: expected %.2f, got %.2f", expectedTaxAmount, req.TaxAmount)
	}
	
	if abs(req.Total-expectedTotal) > tolerance {
		return fmt.Errorf("total mismatch: expected %.2f, got %.2f", expectedTotal, req.Total)
	}
	
	return nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func CreateInvoice(c *gin.Context) {
	var req invoiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Validate date
	date, err := validateDate(req.Date)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Validate calculations
	if err := validateInvoiceCalculations(&req); err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Check for duplicate invoice number for this user
	var existingInv model.Invoice
	if err := db.DB.Where("number = ? AND user_id = ?", req.InvoiceNumber, c.GetUint("userID")).First(&existingInv).Error; err == nil {
		handleError(c, http.StatusConflict, "Invoice number already exists", nil)
		return
	}

	inv := model.Invoice{
		UserID:    c.GetUint("userID"),
		Number:    req.InvoiceNumber,
		Date:      datatypes.Date(date),
		FromName:  req.FromName,
		FromEmail: req.FromEmail,
		ToName:    req.ToName,
		ToEmail:   req.ToEmail,
		TaxRate:   req.TaxRate,
		Subtotal:  req.Subtotal,
		TaxAmount: req.TaxAmount,
		Total:     req.Total,
		Items:     req.Items,
	}

	if err := db.DB.Create(&inv).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create invoice", err)
		return
	}

	c.JSON(http.StatusCreated,gin.H{
		"status": http.StatusCreated,
		"data": inv,
	})
}

func ListInvoices(c *gin.Context) {
	var invoices []model.Invoice
	
	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Search parameter
	search := c.Query("search")
	
	query := db.DB.Preload("Items").Where("user_id = ?", c.GetUint("userID"))
	
	if search != "" {
		query = query.Where("number ILIKE ? OR from_name ILIKE ? OR to_name ILIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Model(&model.Invoice{}).Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&invoices).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to fetch invoices", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       invoices,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func GetInvoice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid invoice ID", err)
		return
	}

	var inv model.Invoice
	if err := db.DB.Preload("Items").
		First(&inv, "id = ? AND user_id = ?", id, c.GetUint("userID")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(c, http.StatusNotFound, "Invoice not found", nil)
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to fetch invoice", err)
		}
		return
	}

	c.JSON(http.StatusOK, inv)
}

func UpdateInvoice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid invoice ID", err)
		return
	}

	var inv model.Invoice
	if err := db.DB.First(&inv, "id = ? AND user_id = ?", id, c.GetUint("userID")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(c, http.StatusNotFound, "Invoice not found", nil)
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to fetch invoice", err)
		}
		return
	}

	var req invoiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Validate date
	date, err := validateDate(req.Date)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Validate calculations
	if err := validateInvoiceCalculations(&req); err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Check for duplicate invoice number (excluding current invoice)
	var existingInv model.Invoice
	if err := db.DB.Where("number = ? AND user_id = ? AND id != ?", req.InvoiceNumber, c.GetUint("userID"), id).First(&existingInv).Error; err == nil {
		handleError(c, http.StatusConflict, "Invoice number already exists", nil)
		return
	}

	// Update invoice fields
	inv.Number = req.InvoiceNumber
	inv.Date = datatypes.Date(date)
	inv.FromName = req.FromName
	inv.FromEmail = req.FromEmail
	inv.ToName = req.ToName
	inv.ToEmail = req.ToEmail
	inv.TaxRate = req.TaxRate
	inv.Subtotal = req.Subtotal
	inv.TaxAmount = req.TaxAmount
	inv.Total = req.Total

	// Use transaction for atomic update
	tx := db.DB.Begin()
	if tx.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to start transaction", tx.Error)
		return
	}

	// Replace items
	if err := tx.Model(&inv).Association("Items").Replace(req.Items); err != nil {
		tx.Rollback()
		handleError(c, http.StatusInternalServerError, "Failed to update invoice items", err)
		return
	}

	// Save invoice
	if err := tx.Save(&inv).Error; err != nil {
		tx.Rollback()
		handleError(c, http.StatusInternalServerError, "Failed to update invoice", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	// Reload with items for response
	db.DB.Preload("Items").First(&inv, inv.ID)
	c.JSON(http.StatusOK, inv)
}

func DeleteInvoice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid invoice ID", err)
		return
	}

	// Use transaction for safe deletion
	tx := db.DB.Begin()
	if tx.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to start transaction", tx.Error)
		return
	}

	res := tx.Where("id = ? AND user_id = ?", id, c.GetUint("userID")).Delete(&model.Invoice{})
	if res.Error != nil {
		tx.Rollback()
		handleError(c, http.StatusInternalServerError, "Failed to delete invoice", res.Error)
		return
	}

	if res.RowsAffected == 0 {
		tx.Rollback()
		handleError(c, http.StatusNotFound, "Invoice not found", nil)
		return
	}

	if err := tx.Commit().Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}