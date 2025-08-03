package restore

import (
	"fmt"
	"hammer/internal"

	"github.com/spf13/cobra"
)

var (
	folderPath string
	connectionString string
	databaseName     string
)

var RestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore your data from a folder",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.RestoreCollections(connectionString, databaseName, folderPath)
	},
}

func init() {
	RestoreCmd.Flags().StringVarP(&folderPath, "folder", "f", "", "Folder path to restore from")
	RestoreCmd.Flags().StringVarP(&connectionString, "connection", "c", "", "Database connection string")
	RestoreCmd.Flags().StringVarP(&databaseName, "database", "d", "", "Database name")
	if err := RestoreCmd.MarkFlagRequired("folder"); err != nil {
		fmt.Println("Folder path not defined")
	}
	if err := RestoreCmd.MarkFlagRequired("connection"); err != nil {
		fmt.Println("Database connection string not defined")
	}
	if err := RestoreCmd.MarkFlagRequired("database"); err != nil {
		fmt.Println("Database name not defined")
	}
}
