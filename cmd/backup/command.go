package backup

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	databaseConnectionString string
)

var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Extract your backup to a zipfile or uploads to aws",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	BackupCmd.Flags().StringVarP(&databaseConnectionString, "db", "d", "", "Database connection string")

	if err := BackupCmd.MarkFlagRequired("db"); err != nil {
		fmt.Println("Database connection string not defined")
	}
}
