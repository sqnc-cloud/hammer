package backup

import (
	"fmt"
	"hammer/internal"

	"github.com/spf13/cobra"
)

var (
	connectionString string
	databaseName     string
	upload           bool
)

var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Extract your backup to a zipfile or uploads to aws",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		backupFile, err := internal.ExportCollections(connectionString, databaseName)
		if err != nil {
			return fmt.Errorf("error to create a backup %v", err)
		}

		if upload {
			awsConfig, err := internal.LoadAWSConfigFromEnv()
			if err != nil {
				return err
			}

			return internal.UploadToS3(backupFile, awsConfig.Bucket, awsConfig.AccessKeyID, awsConfig.SecretAccessKey, awsConfig.Region)
		}
		return nil
	},
}

func init() {
	BackupCmd.Flags().StringVarP(&connectionString, "connection", "c", "", "Database connection string")
	BackupCmd.Flags().StringVarP(&databaseName, "database", "d", "", "Database name")
	BackupCmd.Flags().BoolVarP(&upload, "upload", "u", false, "Upload to aws")

	if err := BackupCmd.MarkFlagRequired("connection"); err != nil {
		fmt.Println("Database connection string not defined")
	}
	if err := BackupCmd.MarkFlagRequired("database"); err != nil {
		fmt.Println("Database name not defined")
	}
}
