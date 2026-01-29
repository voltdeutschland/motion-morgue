package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"motion-morgue/db"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "PDF zu einem Eintrag öffnen",
}

var openAssemblyCmd = &cobra.Command{
	Use:   "assembly [id]",
	Short: "Protokoll-PDF einer Versammlung öffnen",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		var pdfPath *string
		err := db.DB.QueryRow("SELECT protocol_pdf FROM assemblies WHERE id = ?", id).Scan(&pdfPath)
		if err != nil {
			return fmt.Errorf("Versammlung nicht gefunden: %w", err)
		}

		if pdfPath == nil {
			return fmt.Errorf("Kein PDF für diese Versammlung hinterlegt")
		}

		return openFile(*pdfPath)
	},
}

var openMotionCmd = &cobra.Command{
	Use:   "motion [id]",
	Short: "PDF eines Antrags öffnen",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		var pdfPath *string
		err := db.DB.QueryRow("SELECT pdf_path FROM motions WHERE id = ?", id).Scan(&pdfPath)
		if err != nil {
			return fmt.Errorf("Antrag nicht gefunden: %w", err)
		}

		if pdfPath == nil {
			return fmt.Errorf("Kein PDF für diesen Antrag hinterlegt")
		}

		return openFile(*pdfPath)
	},
}

var openAmendmentCmd = &cobra.Command{
	Use:   "amendment [id]",
	Short: "PDF eines Änderungsantrags öffnen",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		var pdfPath *string
		err := db.DB.QueryRow("SELECT pdf_path FROM amendments WHERE id = ?", id).Scan(&pdfPath)
		if err != nil {
			return fmt.Errorf("Änderungsantrag nicht gefunden: %w", err)
		}

		if pdfPath == nil {
			return fmt.Errorf("Kein PDF für diesen Änderungsantrag hinterlegt")
		}

		return openFile(*pdfPath)
	},
}

func openFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		return fmt.Errorf("Betriebssystem nicht unterstützt: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.AddCommand(openAssemblyCmd)
	openCmd.AddCommand(openMotionCmd)
	openCmd.AddCommand(openAmendmentCmd)
}
