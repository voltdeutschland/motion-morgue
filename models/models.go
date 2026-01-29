package models

import "database/sql"

type Assembly struct {
	ID          int64
	Title       string
	StartDate   sql.NullString
	EndDate     sql.NullString
	ProtocolPDF sql.NullString
}

type Motion struct {
	ID         int64
	AssemblyID int64
	Title      string
	SortNumber string
	PDFPath    sql.NullString
}

type Amendment struct {
	ID         int64
	MotionID   int64
	SortNumber string
	PDFPath    sql.NullString
}
