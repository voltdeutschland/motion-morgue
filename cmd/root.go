package cmd

import (
	"fmt"
	"os"

	"motion-morgue/db"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "motion-morgue",
	Short: "Verwaltung von Versammlungsdokumenten",
	Long:  `Ein CLI-Tool zur Verwaltung von Versammlungsprotokollen, Anträgen und Änderungsanträgen.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return db.Init()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.Close()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
