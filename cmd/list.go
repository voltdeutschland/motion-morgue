package cmd

import (
	"fmt"
	"strings"

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
	rows, err := db.DB.Query("SELECT id, title, start_date, end_date, protocol_pdf FROM assemblies ORDER BY start_date DESC")
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
	rows, err := db.DB.Query("SELECT id, assembly_id, title, sort_number, pdf_path FROM motions WHERE assembly_id = ? ORDER BY sort_number", assemblyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var motions []models.Motion
	for rows.Next() {
		var m models.Motion
		if err := rows.Scan(&m.ID, &m.AssemblyID, &m.Title, &m.SortNumber, &m.PDFPath); err != nil {
			return nil, err
		}
		motions = append(motions, m)
	}
	return motions, nil
}

func getAmendments(motionID int64) ([]models.Amendment, error) {
	rows, err := db.DB.Query("SELECT id, motion_id, title, sort_number, pdf_path FROM amendments WHERE motion_id = ? ORDER BY sort_number", motionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var amendments []models.Amendment
	for rows.Next() {
		var a models.Amendment
		if err := rows.Scan(&a.ID, &a.MotionID, &a.Title, &a.SortNumber, &a.PDFPath); err != nil {
			return nil, err
		}
		amendments = append(amendments, a)
	}
	return amendments, nil
}

func printAssembly(a models.Assembly) {
	dateRange := ""
	if a.StartDate.Valid {
		dateRange = a.StartDate.String
		if a.EndDate.Valid {
			dateRange += " - " + a.EndDate.String
		}
	}

	pdfMarker := "[ - ]"
	if a.ProtocolPDF.Valid {
		pdfMarker = "[PDF]"
	}

	header := fmt.Sprintf(" %s", a.Title)
	if dateRange != "" {
		header += fmt.Sprintf(" (%s)", dateRange)
	}

	width := 70
	headerPadding := width - len(header) - len(pdfMarker) - 3
	if headerPadding < 1 {
		headerPadding = 1
	}

	fmt.Println("┌" + strings.Repeat("─", width) + "┐")
	fmt.Printf("│%s%s%s │\n", header, strings.Repeat(" ", headerPadding), pdfMarker)
	fmt.Println("├" + strings.Repeat("─", width) + "┤")

	motions, _ := getMotions(a.ID)
	for _, m := range motions {
		motionPDF := "[ - ]"
		if m.PDFPath.Valid {
			motionPDF = "[PDF]"
		}

		motionLine := fmt.Sprintf("   %s  %s", m.SortNumber, m.Title)
		motionPadding := width - len(motionLine) - len(motionPDF) - 3
		if motionPadding < 1 {
			motionPadding = 1
		}
		fmt.Printf("│%s%s%s │\n", motionLine, strings.Repeat(" ", motionPadding), motionPDF)

		amendments, _ := getAmendments(m.ID)
		for _, am := range amendments {
			amendPDF := "[ - ]"
			if am.PDFPath.Valid {
				amendPDF = "[PDF]"
			}

			amendLine := fmt.Sprintf("     └─ %s", am.SortNumber)
			if am.Title.Valid {
				amendLine += "  " + am.Title.String
			}
			amendPadding := width - len(amendLine) - len(amendPDF) - 3
			if amendPadding < 1 {
				amendPadding = 1
			}
			fmt.Printf("│%s%s%s │\n", amendLine, strings.Repeat(" ", amendPadding), amendPDF)
		}
	}

	fmt.Println("└" + strings.Repeat("─", width) + "┘")
}

func init() {
	rootCmd.AddCommand(listCmd)
}
