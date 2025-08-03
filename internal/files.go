package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func WriteBSONCollectionToFile(cursor *mongo.Cursor, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	for cursor.Next(context.TODO()) {
		_, err := file.Write(cursor.Current)
		if err != nil {
			return fmt.Errorf("failed to write BSON data to file %s: %w", filePath, err)
		}
	}

	return nil
}

func WriteJSONCollectionToFile(cursor *mongo.Cursor, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	for cursor.Next(context.TODO()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return fmt.Errorf("failed to decode BSON to JSON: %w", err)
		}
		if err := encoder.Encode(doc); err != nil {
			return fmt.Errorf("failed to write JSON data to file %s: %w", filePath, err)
		}
	}

	return nil
}
