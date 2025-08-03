package root

import (
	"github.com/spf13/cobra"
	"hammer/cmd/backup"
	"hammer/cmd/restore"
	"os"
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
	rootCmd.AddCommand(restore.RestoreCmd)
}

func init() {
	addSubcommandPalettes()
}
