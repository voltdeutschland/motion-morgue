package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"motion-morgue/db"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "PDF importieren und verknüpfen",
}

var importAssemblyCmd = &cobra.Command{
	Use:   "assembly [id] [pdf-path]",
	Short: "Protokoll-PDF zu Versammlung importieren",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		srcPath := args[1]

		destPath, err := copyPDF(srcPath, "assembly", id)
		if err != nil {
			return err
		}

		if _, err := db.DB.Exec("UPDATE assemblies SET protocol_pdf = ? WHERE id = ?", destPath, id); err != nil {
			return err
		}

		fmt.Printf("PDF importiert: %s\n", destPath)
		return nil
	},
}

var importMotionCmd = &cobra.Command{
	Use:   "motion [id] [pdf-path]",
	Short: "PDF zu Antrag importieren",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		srcPath := args[1]

		destPath, err := copyPDF(srcPath, "motion", id)
		if err != nil {
			return err
		}

		if _, err := db.DB.Exec("UPDATE motions SET pdf_path = ? WHERE id = ?", destPath, id); err != nil {
			return err
		}

		fmt.Printf("PDF importiert: %s\n", destPath)
		return nil
	},
}

var importAmendmentCmd = &cobra.Command{
	Use:   "amendment [id] [pdf-path]",
	Short: "PDF zu Änderungsantrag importieren",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		srcPath := args[1]

		destPath, err := copyPDF(srcPath, "amendment", id)
		if err != nil {
			return err
		}

		if _, err := db.DB.Exec("UPDATE amendments SET pdf_path = ? WHERE id = ?", destPath, id); err != nil {
			return err
		}

		fmt.Printf("PDF importiert: %s\n", destPath)
		return nil
	},
}

func copyPDF(srcPath, entityType, id string) (string, error) {
	dataDir, err := db.DataDir()
	if err != nil {
		return "", err
	}

	pdfDir := filepath.Join(dataDir, "pdfs")
	if err := os.MkdirAll(pdfDir, 0755); err != nil {
		return "", fmt.Errorf("could not create pdf directory: %w", err)
	}

	filename := fmt.Sprintf("%s_%s%s", entityType, id, filepath.Ext(srcPath))
	destPath := filepath.Join(pdfDir, filename)

	src, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("could not open source file: %w", err)
	}
	defer src.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("could not create destination file: %w", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return "", fmt.Errorf("could not copy file: %w", err)
	}

	return destPath, nil
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.AddCommand(importAssemblyCmd)
	importCmd.AddCommand(importMotionCmd)
	importCmd.AddCommand(importAmendmentCmd)
}
