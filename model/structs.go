package model

import "time"

// ===================== Request Structs ======================================================= //

// RequestInput struct
type RequestInput struct {
	Date     string `json:"date"`
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

// RequestDates struct
type RequestDates struct {
	DateFrom time.Time
	DateTo   time.Time
}
