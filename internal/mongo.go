package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	return client, nil
}

func DisconnectMongoDB(client *mongo.Client) {
	if client == nil {
		return
	}
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(fmt.Errorf("failed to disconnect from MongoDB: %w", err))
	}
}

func ListMongoDBCollectionNames(db *mongo.Database) ([]string, error) {
	collections, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collection names: %w", err)
	}
	return collections, nil
}

func ExportCollections(uri, databaseName string) (string, error) {
	client, err := ConnectMongoDB(uri)
	if err != nil {
		return "", err
	}
	defer DisconnectMongoDB(client)

	db := client.Database(databaseName)
	collections, err := ListMongoDBCollectionNames(db)
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "mongo-export-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	for _, collection := range collections {
		fmt.Printf("Exporting collection: %s\n", collection)

		cursor, err := db.Collection(collection).Find(context.TODO(), bson.D{})
		if err != nil {
			return "", fmt.Errorf("failed to find documents in collection %s: %w", collection, err)
		}
		defer cursor.Close(context.TODO())

		filePath := filepath.Join(tempDir, fmt.Sprintf("%s.json", collection))
		if err := WriteJSONCollectionToFile(cursor, filePath); err != nil {
			return "", fmt.Errorf("failed to write collection %s to file: %w", collection, err)
		}
	}

	outputFile := "backup.zip"
	if err := ZipFolder(tempDir, outputFile); err != nil {
		return "", fmt.Errorf("failed to zip collection files: %w", err)
	}

	return outputFile, nil
}

func RestoreCollections(uri, databaseName, folderPath string) error {
	client, err := ConnectMongoDB(uri)
	if err != nil {
		return err
	}
	defer DisconnectMongoDB(client)

	db := client.Database(databaseName)

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			collectionName := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			fmt.Printf("Importing collection: %s\n", collectionName)

			filePath := filepath.Join(folderPath, file.Name())
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", filePath, err)
			}

			var docs []interface{}
			if err := bson.UnmarshalExtJSON(fileData, true, &docs); err != nil {
				return fmt.Errorf("failed to unmarshal json for file %s: %w", filePath, err)
			}

			collection := db.Collection(collectionName)
			_, err = collection.InsertMany(context.TODO(), docs)
			if err != nil {
				return fmt.Errorf("failed to insert documents into collection %s: %w", collectionName, err)
			}
		}
	}

	return nil
}
