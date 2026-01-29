package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"motion-morgue/db"

	"github.com/spf13/cobra"
)

var reader = bufio.NewReader(os.Stdin)

func prompt(label string, current string) string {
	if current != "" {
		return current
	}
	fmt.Printf("%s: ", label)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptInt(label string, current int64) int64 {
	if current != 0 {
		return current
	}
	fmt.Printf("%s: ", label)
	input, _ := reader.ReadString('\n')
	val, _ := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
	return val
}

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

		title = prompt("Titel", title)
		start = prompt("Startdatum (YYYY-MM-DD, optional)", start)
		end = prompt("Enddatum (YYYY-MM-DD, optional)", end)

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
		status, _ := cmd.Flags().GetString("status")

		assemblyID = promptInt("Assembly ID", assemblyID)
		sort = prompt("Sortiernummer (z.B. A001)", sort)
		title = prompt("Titel", title)
		status = prompt("Status (d|s|w|a|b|p|r, optional)", status)

		result, err := db.DB.Exec(
			"INSERT INTO motions (assembly_id, sort_number, title, status) VALUES (?, ?, ?, ?)",
			assemblyID, sort, title, nullIfEmpty(status),
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
		status, _ := cmd.Flags().GetString("status")

		motionID = promptInt("Motion ID", motionID)
		sort = prompt("Sortiernummer (z.B. Ä001)", sort)
		title = prompt("Titel (optional)", title)
		status = prompt("Status (d|s|w|a|b|p|r|m, optional)", status)

		result, err := db.DB.Exec(
			"INSERT INTO amendments (motion_id, sort_number, title, status) VALUES (?, ?, ?, ?)",
			motionID, sort, nullIfEmpty(title), nullIfEmpty(status),
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

	addCmd.AddCommand(addMotionCmd)
	addMotionCmd.Flags().Int64("assembly", 0, "ID der Versammlung")
	addMotionCmd.Flags().String("sort", "", "Sortiernummer (z.B. A001)")
	addMotionCmd.Flags().String("title", "", "Titel des Antrags")
	addMotionCmd.Flags().String("status", "", "Status (d|s|w|a|b|p|r)")

	addCmd.AddCommand(addAmendmentCmd)
	addAmendmentCmd.Flags().Int64("motion", 0, "ID des Antrags")
	addAmendmentCmd.Flags().String("sort", "", "Sortiernummer (z.B. Ä001)")
	addAmendmentCmd.Flags().String("title", "", "Titel des Änderungsantrags")
	addAmendmentCmd.Flags().String("status", "", "Status (d|s|w|a|b|p|r|m)")
}
