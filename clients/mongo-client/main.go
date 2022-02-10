package main

import (
	"fmt"
	"os"

	mongo "e.co/m/db"
	"e.co/m/vault"
	"github.com/fatih/color"
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

	client := vault.DynamicSecrets{
		Token: token,
	}
	creds := client.GetCredentials()
	repo := mongo.Repo{}

	message := &mongo.Message{
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
