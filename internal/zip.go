package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// ZipFolder zips the content of the source folder to the destination zip file.
func ZipFolder(sourceFolder, destZipFile string) error {
	zipFile, err := os.Create(destZipFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	// Get the base directory name to use as the root within the zip file
	// If sourceFolder is "/tmp/mydata", baseDir will be "mydata"
	baseDir := filepath.Base(sourceFolder)

	filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a header for the file/directory
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Determine the name within the zip file
		// This ensures that the zip file contains the base directory as its root
		relPath, err := filepath.Rel(filepath.Dir(sourceFolder), path)
		if err != nil {
			return err
		}
		header.Name = filepath.Join(baseDir, relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}
		return nil
	})

	return nil
}

