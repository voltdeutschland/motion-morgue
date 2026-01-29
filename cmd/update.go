package cmd

import (
	"fmt"

	"motion-morgue/db"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Einträge aktualisieren",
}

var updateAssemblyCmd = &cobra.Command{
	Use:   "assembly [id]",
	Short: "Versammlung aktualisieren",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		title, _ := cmd.Flags().GetString("title")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		pdf, _ := cmd.Flags().GetString("pdf")

		if title != "" {
			if _, err := db.DB.Exec("UPDATE assemblies SET title = ? WHERE id = ?", title, id); err != nil {
				return err
			}
		}
		if start != "" {
			if _, err := db.DB.Exec("UPDATE assemblies SET start_date = ? WHERE id = ?", start, id); err != nil {
				return err
			}
		}
		if end != "" {
			if _, err := db.DB.Exec("UPDATE assemblies SET end_date = ? WHERE id = ?", end, id); err != nil {
				return err
			}
		}
		if pdf != "" {
			if _, err := db.DB.Exec("UPDATE assemblies SET protocol_pdf = ? WHERE id = ?", pdf, id); err != nil {
				return err
			}
		}

		fmt.Println("Versammlung aktualisiert")
		return nil
	},
}

var updateMotionCmd = &cobra.Command{
	Use:   "motion [id]",
	Short: "Antrag aktualisieren",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		title, _ := cmd.Flags().GetString("title")
		sort, _ := cmd.Flags().GetString("sort")
		status, _ := cmd.Flags().GetString("status")
		pdf, _ := cmd.Flags().GetString("pdf")

		if title != "" {
			if _, err := db.DB.Exec("UPDATE motions SET title = ? WHERE id = ?", title, id); err != nil {
				return err
			}
		}
		if sort != "" {
			if _, err := db.DB.Exec("UPDATE motions SET sort_number = ? WHERE id = ?", sort, id); err != nil {
				return err
			}
		}
		if status != "" {
			if _, err := db.DB.Exec("UPDATE motions SET status = ? WHERE id = ?", status, id); err != nil {
				return err
			}
		}
		if pdf != "" {
			if _, err := db.DB.Exec("UPDATE motions SET pdf_path = ? WHERE id = ?", pdf, id); err != nil {
				return err
			}
		}

		fmt.Println("Antrag aktualisiert")
		return nil
	},
}

var updateAmendmentCmd = &cobra.Command{
	Use:   "amendment [id]",
	Short: "Änderungsantrag aktualisieren",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		title, _ := cmd.Flags().GetString("title")
		sort, _ := cmd.Flags().GetString("sort")
		status, _ := cmd.Flags().GetString("status")
		pdf, _ := cmd.Flags().GetString("pdf")

		if title != "" {
			if _, err := db.DB.Exec("UPDATE amendments SET title = ? WHERE id = ?", title, id); err != nil {
				return err
			}
		}
		if sort != "" {
			if _, err := db.DB.Exec("UPDATE amendments SET sort_number = ? WHERE id = ?", sort, id); err != nil {
				return err
			}
		}
		if status != "" {
			if _, err := db.DB.Exec("UPDATE amendments SET status = ? WHERE id = ?", status, id); err != nil {
				return err
			}
		}
		if pdf != "" {
			if _, err := db.DB.Exec("UPDATE amendments SET pdf_path = ? WHERE id = ?", pdf, id); err != nil {
				return err
			}
		}

		fmt.Println("Änderungsantrag aktualisiert")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.AddCommand(updateAssemblyCmd)
	updateAssemblyCmd.Flags().String("title", "", "Neuer Titel")
	updateAssemblyCmd.Flags().String("start", "", "Neues Startdatum")
	updateAssemblyCmd.Flags().String("end", "", "Neues Enddatum")
	updateAssemblyCmd.Flags().String("pdf", "", "Pfad zum Protokoll-PDF")

	updateCmd.AddCommand(updateMotionCmd)
	updateMotionCmd.Flags().String("title", "", "Neuer Titel")
	updateMotionCmd.Flags().String("sort", "", "Neue Sortiernummer")
	updateMotionCmd.Flags().String("status", "", "Neuer Status")
	updateMotionCmd.Flags().String("pdf", "", "Pfad zum PDF")

	updateCmd.AddCommand(updateAmendmentCmd)
	updateAmendmentCmd.Flags().String("title", "", "Neuer Titel")
	updateAmendmentCmd.Flags().String("sort", "", "Neue Sortiernummer")
	updateAmendmentCmd.Flags().String("status", "", "Neuer Status")
	updateAmendmentCmd.Flags().String("pdf", "", "Pfad zum PDF")
}
