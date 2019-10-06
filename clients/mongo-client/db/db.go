package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Message - Test Struct
type Message struct {
	Subject string
	Content string
}

// Repo - Database Operations
type Repo struct{}

// Connect - Connect to MongoDB test db and return my-collection
func (repo Repo) Connect(username string, password string) *mongo.Client {
	mongoClientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@localhost:27017", username, password))

	mongoClient, err := mongo.Connect(context.TODO(), mongoClientOptions)

	if err != nil {
		fmt.Println("trying to create mongo client...")
		log.Fatal(err)
	}

	// Check the connection
	if err := mongoClient.Ping(context.TODO(), nil); err != nil {
		fmt.Println("checking mongoclient connection...")
		log.Fatal(err)
	}

	// var result TestMongoResult

	return mongoClient

}

// Insert - Insert data into my-collection
func (repo Repo) Insert(client *mongo.Client, db string, collectionName string, message *Message) {

	collection := client.Database(db).Collection(collectionName)

	insertResult, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		color.Red(fmt.Sprintf("\nMongoDB Request: 403 [Forbidden]\n\t%s\n", err))
		fmt.Println("")
		os.Exit(2)
	}

	color.HiGreen(fmt.Sprintf("\nMongoDB Request: 200 [OK]\nInserted a single document: %v into %s database\n", insertResult.InsertedID, db))
	fmt.Println("")

}

// Read - Read my-collection and return the data
func (repo Repo) Read(client *mongo.Client, db string, collectionName string, message *Message) {

}

// Update - Update a record in my-collection
func (repo Repo) Update(client *mongo.Client, db string, collectionName string, message *Message) {

}

// Delete - Delete a record in my-collection
func (repo Repo) Delete(client *mongo.Client, db string, collectionName string, message *Message) {

}
