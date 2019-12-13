package model

import "time"

// ReportRequest struct
type ReportRequest struct {
	Dates      *RequestDates
	ReportType *ReportType
}

// RequestInput struct
type RequestInput struct {
	Date       string `json:"date"`
	DateFrom   string `json:"dateFrom"`
	DateTo     string `json:"dateTo"`
	ReportType string `json:"type"`
}

// RequestDates struct
type RequestDates struct {
	DateFrom time.Time
	DateTo   time.Time
}
