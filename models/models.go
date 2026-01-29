package models

import "database/sql"

type Assembly struct {
	ID          int64
	Title       string
	StartDate   sql.NullString
	EndDate     sql.NullString
	ProtocolPDF sql.NullString
}

// Valid motion statuses: draft, submitted, withdrawn, admitted, not_admitted, approved, rejected
type Motion struct {
	ID         int64
	AssemblyID int64
	Title      string
	SortNumber string
	Status     sql.NullString
	PDFPath    sql.NullString
}

// Valid amendment statuses: draft, submitted, withdrawn, admitted, not_admitted, approved, rejected, adopted
type Amendment struct {
	ID         int64
	MotionID   int64
	Title      sql.NullString
	SortNumber string
	Status     sql.NullString
	PDFPath    sql.NullString
}
