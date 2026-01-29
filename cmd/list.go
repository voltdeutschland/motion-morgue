package cmd

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"motion-morgue/db"
	"motion-morgue/models"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Hierarchische Ansicht aller Versammlungen",
	RunE: func(cmd *cobra.Command, args []string) error {
		assemblies, err := getAssemblies()
		if err != nil {
			return err
		}

		for _, a := range assemblies {
			printAssembly(a)
		}

		return nil
	},
}

func getAssemblies() ([]models.Assembly, error) {
	rows, err := db.DB.Query("SELECT id, title, start_date, end_date, protocol_pdf FROM assemblies ORDER BY start_date ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assemblies []models.Assembly
	for rows.Next() {
		var a models.Assembly
		if err := rows.Scan(&a.ID, &a.Title, &a.StartDate, &a.EndDate, &a.ProtocolPDF); err != nil {
			return nil, err
		}
		assemblies = append(assemblies, a)
	}
	return assemblies, nil
}

func getMotions(assemblyID int64) ([]models.Motion, error) {
	rows, err := db.DB.Query("SELECT id, assembly_id, title, sort_number, status, pdf_path FROM motions WHERE assembly_id = ? ORDER BY sort_number", assemblyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var motions []models.Motion
	for rows.Next() {
		var m models.Motion
		if err := rows.Scan(&m.ID, &m.AssemblyID, &m.Title, &m.SortNumber, &m.Status, &m.PDFPath); err != nil {
			return nil, err
		}
		motions = append(motions, m)
	}
	return motions, nil
}

func getAmendments(motionID int64) ([]models.Amendment, error) {
	rows, err := db.DB.Query("SELECT id, motion_id, title, sort_number, status, pdf_path FROM amendments WHERE motion_id = ? ORDER BY sort_number", motionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var amendments []models.Amendment
	for rows.Next() {
		var a models.Amendment
		if err := rows.Scan(&a.ID, &a.MotionID, &a.Title, &a.SortNumber, &a.Status, &a.PDFPath); err != nil {
			return nil, err
		}
		amendments = append(amendments, a)
	}
	return amendments, nil
}

func statusMarker(status string) string {
	switch status {
	case "d", "draft":
		return " D "
	case "s", "submitted":
		return " S "
	case "w", "withdrawn":
		return " W "
	case "a", "admitted":
		return " A "
	case "b", "blocked":
		return " B "
	case "p", "passed":
		return " ✓ "
	case "r", "rejected":
		return " ✗ "
	case "m", "merged":
		return " M "
	default:
		return " ? "
	}
}

const tableWidth = 80

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

func truncate(s string, maxLen int) string {
	if runeLen(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string([]rune(s)[:maxLen])
	}
	return string([]rune(s)[:maxLen-3]) + "..."
}

func printLine(content string, status string, pdfMarker string) {
	marker := fmt.Sprintf("[%s|%s]", status, pdfMarker)
	maxContent := tableWidth - runeLen(marker) - 6
	content = truncate(content, maxContent)
	padding := tableWidth - runeLen(content) - runeLen(marker) - 3
	fmt.Printf("│%s%s%s   │\n", content, strings.Repeat(" ", padding), marker)
}

func printAssembly(a models.Assembly) {
	dateRange := ""
	if a.StartDate.Valid {
		dateRange = a.StartDate.String
		if a.EndDate.Valid {
			dateRange += " - " + a.EndDate.String
		}
	}

	pdfMarker := " - "
	if a.ProtocolPDF.Valid {
		pdfMarker = "PDF"
	}

	header := fmt.Sprintf(" %s", a.Title)
	if dateRange != "" {
		header += fmt.Sprintf(" (%s)", dateRange)
	}

	fmt.Println("┌" + strings.Repeat("─", tableWidth) + "┐")
	printLine(fmt.Sprintf(" #%d %s", a.ID, header[1:]), " - ", pdfMarker)
	fmt.Println("├" + strings.Repeat("─", tableWidth) + "┤")

	motions, _ := getMotions(a.ID)
	for _, m := range motions {
		motionPDF := " - "
		if m.PDFPath.Valid {
			motionPDF = "PDF"
		}
		motionStatus := statusMarker(m.Status.String)
		if !m.Status.Valid {
			motionStatus = " ? "
		}
		printLine(fmt.Sprintf("   #%-3d %s  %s", m.ID, m.SortNumber, m.Title), motionStatus, motionPDF)

		amendments, _ := getAmendments(m.ID)
		for _, am := range amendments {
			amendPDF := " - "
			if am.PDFPath.Valid {
				amendPDF = "PDF"
			}
			amendStatus := statusMarker(am.Status.String)
			if !am.Status.Valid {
				amendStatus = " ? "
			}
			amendLine := fmt.Sprintf("        #%-3d └─ %s", am.ID, am.SortNumber)
			if am.Title.Valid {
				amendLine += "  " + am.Title.String
			}
			printLine(amendLine, amendStatus, amendPDF)
		}
	}

	fmt.Println("└" + strings.Repeat("─", tableWidth) + "┘")
}

func init() {
	rootCmd.AddCommand(listCmd)
}
