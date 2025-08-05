package model

import (
	"time"

	"gorm.io/datatypes"
)

type Invoice struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	UserID       uint          `json:"userID"`
	Number       string        `gorm:"uniqueIndex;size:50" json:"invoiceNumber"`
	Date         datatypes.Date `json:"date"`
	FromName     string        `json:"fromName"`
	FromEmail    string        `json:"fromEmail"`
	ToName       string        `json:"toName"`
	ToEmail      string        `json:"toEmail"`
	TaxRate      float64       `json:"taxRate"`
	Subtotal     float64       `json:"subtotal"`
	TaxAmount    float64       `json:"taxAmount"`
	Total        float64       `json:"total"`
	Items        []InvoiceItem `gorm:"foreignKey:InvoiceID" json:"items"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

type InvoiceItem struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	InvoiceID   uint    `json:"-"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Rate        float64 `json:"rate"`
	Amount      float64 `json:"amount"`
}