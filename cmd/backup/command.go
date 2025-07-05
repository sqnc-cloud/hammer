package backup

import (
	"fmt"
	"github.com/spf13/cobra"
	"hammer/internal"
)

var (
	connectionString string
)

var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Extract your backup to a zipfile or uploads to aws",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		uri := "mongodb://localhost:27017"
		database := "hammer"
		internal.ExportCollections(uri, database)
	},
}

func init() {
	BackupCmd.Flags().StringVarP(&connectionString, "db", "d", "", "Database connection string")

	if err := BackupCmd.MarkFlagRequired("db"); err != nil {
		fmt.Println("Database connection string not defined")
	}
}
