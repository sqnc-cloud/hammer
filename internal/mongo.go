package internal

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"os/exec"
)

func ExportCollections(uri string, database string) {
	collections := collectDatabaseCollections(database, uri)
	fmt.Println("Found:", collections)
	for _, collection := range collections {
		exportCollection(collection, uri)
	}
}

func exportCollection(collection string, uri string) {
	fmt.Println("Exporting collection ", collection)
	cmd := exec.Command(
		"mongoexport",
		"--uri", uri,
		"--db", "hammer",
		"--collection", collection,
		"--out", fmt.Sprintf("./backup/%s.json", collection),
		"--jsonFormat", "canonical",
		"--type", "json")

	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func collectDatabaseCollections(database string, uri string) []string {
	fmt.Println("Listing collections")
	client := connectDatabase(uri)
	db := client.Database(database)
	result, err := db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func connectDatabase(uri string) *mongo.Client {
	fmt.Println("Connecting to database.")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	return client
}
