package cmd

import (
	"hammer/cmd/backup"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hammer",
	Short: "A tool to backup your mongodb collections into aws s3",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommandPalettes() {
	rootCmd.AddCommand(backup.BackupCmd)
}

func init() {

	addSubcommandPalettes()
}
