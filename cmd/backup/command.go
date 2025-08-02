package backup

import (
	"fmt"
	"github.com/spf13/cobra"
	"hammer/internal"
)

var (
	connectionString string
	upload             bool
)

var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Extract your backup to a zipfile or uploads to aws",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		uri := "mongodb://localhost:27017"
		database := "hammer"
		backupFile, err := internal.ExportCollections(uri, database)

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
	BackupCmd.Flags().StringVarP(&connectionString, "db", "d", "", "Database connection string")
	BackupCmd.Flags().BoolVarP(&upload, "upload", "u", false, "Upload to aws")

	if err := BackupCmd.MarkFlagRequired("db"); err != nil {
		fmt.Println("Database connection string not defined")
	}
}
