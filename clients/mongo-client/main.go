package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"context"
	"log"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// Token: "s.lPV34WdnBvKqxVuCkPLJ2oNN",
	var token string
	var operation string
	var db string

	if len(os.Args) >= 3 {
		token = os.Args[1]
		operation = os.Args[2]
		db = os.Args[3]
	} else {
		color.Red("Please supply your client Vault token and the DB you wish to test ex. (dev, test, production)")
		os.Exit(1)
	}

	fmt.Println(token)

	client := DynamicSecrets{
		Token: token,
	}
	creds := client.GetCredentials()
	repo := Repo{}

	message := &Message{
		Subject: fmt.Sprintf("%s DB Message", db),
		Content: fmt.Sprintf("%s Database Write Success!", db),
	}

	mongoClient := repo.Connect(creds.Username, creds.Password)

	switch operation {
	case "insert":
		repo.Insert(mongoClient, db, "messages", message)
	case "read":
		repo.Read(mongoClient, db, "messages", message)
	case "update":
		repo.Read(mongoClient, db, "messages", message)
	case "delete":
		repo.Read(mongoClient, db, "messages", message)

	}

}

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

	color.Green(fmt.Sprintf("\nMongoDB Request: 200 [OK]\nInserted a single document: %v into %s database\n", insertResult.InsertedID, db))
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

// DynamicSecrets -
type DynamicSecrets struct {
	Token string
}

// Credential -
type Credential struct {
	Username string
	Password string
}

// GetCredentials - return dynamic credentials
func (vault DynamicSecrets) GetCredentials() *Credential {
	var VaultJSONResponse map[string]interface{}
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://127.0.0.1:8200/v1/database/creds/mongo-db-admin", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("X-Vault-Token", vault.Token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal(bodyBytes, &VaultJSONResponse); err != nil {
			log.Fatal(err)
		}

		creds := VaultJSONResponse["data"].(map[string]interface{})
		username := creds["username"].(string)
		password := creds["password"].(string)
		lease := VaultJSONResponse["lease_duration"]
		renewable := VaultJSONResponse["renewable"]
		leaseID := VaultJSONResponse["lease_id"]

		color.Cyan(fmt.Sprintf("\nVault Request: %v [OK]  You have successfully authenticated with Vault and a MongoDB credential has been created and will last for %v seconds\n", resp.StatusCode, VaultJSONResponse["lease_duration"]))
		fmt.Println("")
		color.Cyan(fmt.Sprintf("\tUserName: %v\n\tPassword: %v\n\tLeaseDuration: %v seconds\n\tRenewable: %v\n\tLeaseID: %v", username, password, lease, renewable, leaseID))
		return &Credential{
			Username: username,
			Password: password,
		}

	} else if resp.StatusCode == http.StatusForbidden {
		color.Red(fmt.Sprintf("\nStatusCode: %v [Forbidden]\nVault cannot identify your application. Your client token is invalid..", resp.StatusCode))
		fmt.Println("")
		os.Exit(1)
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("")
		color.Red(fmt.Sprintf("\nStatusCode: %v [Bad Request]\nPlease provide your Vault authorization token..\n", resp.StatusCode))
		os.Exit(1)
	}

	panic(fmt.Sprintf("unexpected error code: [%v]", resp.StatusCode))
}
