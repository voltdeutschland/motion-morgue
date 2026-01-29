package cmd

import (
	"fmt"

	"motion-morgue/db"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Neue Einträge anlegen",
}

var addAssemblyCmd = &cobra.Command{
	Use:   "assembly",
	Short: "Neue Versammlung anlegen",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")

		result, err := db.DB.Exec(
			"INSERT INTO assemblies (title, start_date, end_date) VALUES (?, ?, ?)",
			title, nullIfEmpty(start), nullIfEmpty(end),
		)
		if err != nil {
			return fmt.Errorf("could not create assembly: %w", err)
		}

		id, _ := result.LastInsertId()
		fmt.Printf("Versammlung angelegt mit ID %d\n", id)
		return nil
	},
}

var addMotionCmd = &cobra.Command{
	Use:   "motion",
	Short: "Neuen Antrag anlegen",
	RunE: func(cmd *cobra.Command, args []string) error {
		assemblyID, _ := cmd.Flags().GetInt64("assembly")
		sort, _ := cmd.Flags().GetString("sort")
		title, _ := cmd.Flags().GetString("title")

		result, err := db.DB.Exec(
			"INSERT INTO motions (assembly_id, sort_number, title) VALUES (?, ?, ?)",
			assemblyID, sort, title,
		)
		if err != nil {
			return fmt.Errorf("could not create motion: %w", err)
		}

		id, _ := result.LastInsertId()
		fmt.Printf("Antrag angelegt mit ID %d\n", id)
		return nil
	},
}

var addAmendmentCmd = &cobra.Command{
	Use:   "amendment",
	Short: "Neuen Änderungsantrag anlegen",
	RunE: func(cmd *cobra.Command, args []string) error {
		motionID, _ := cmd.Flags().GetInt64("motion")
		sort, _ := cmd.Flags().GetString("sort")
		title, _ := cmd.Flags().GetString("title")

		result, err := db.DB.Exec(
			"INSERT INTO amendments (motion_id, sort_number, title) VALUES (?, ?, ?)",
			motionID, sort, nullIfEmpty(title),
		)
		if err != nil {
			return fmt.Errorf("could not create amendment: %w", err)
		}

		id, _ := result.LastInsertId()
		fmt.Printf("Änderungsantrag angelegt mit ID %d\n", id)
		return nil
	},
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.AddCommand(addAssemblyCmd)
	addAssemblyCmd.Flags().String("title", "", "Titel der Versammlung")
	addAssemblyCmd.Flags().String("start", "", "Startdatum (YYYY-MM-DD)")
	addAssemblyCmd.Flags().String("end", "", "Enddatum (YYYY-MM-DD)")
	addAssemblyCmd.MarkFlagRequired("title")

	addCmd.AddCommand(addMotionCmd)
	addMotionCmd.Flags().Int64("assembly", 0, "ID der Versammlung")
	addMotionCmd.Flags().String("sort", "", "Sortiernummer (z.B. A001)")
	addMotionCmd.Flags().String("title", "", "Titel des Antrags")
	addMotionCmd.MarkFlagRequired("assembly")
	addMotionCmd.MarkFlagRequired("sort")
	addMotionCmd.MarkFlagRequired("title")

	addCmd.AddCommand(addAmendmentCmd)
	addAmendmentCmd.Flags().Int64("motion", 0, "ID des Antrags")
	addAmendmentCmd.Flags().String("sort", "", "Sortiernummer (z.B. Ä001)")
	addAmendmentCmd.Flags().String("title", "", "Titel des Änderungsantrags (optional)")
	addAmendmentCmd.MarkFlagRequired("motion")
	addAmendmentCmd.MarkFlagRequired("sort")
}
