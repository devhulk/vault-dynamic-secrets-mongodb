package main

import (
	"fmt"
	"os"

	mongo "./db"
	"./vault"
	"github.com/fatih/color"
)

func main() {

	// Token: "s.lPV34WdnBvKqxVuCkPLJ2oNN",
	var token string
	var db string

	if len(os.Args) >= 3 {
		token = os.Args[1]
		db = os.Args[2]
	} else {
		color.Red("Please supply your client Vault token and the DB you wish to test ex. (dev, test, production)")
		os.Exit(1)
	}

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

	repo.Insert(mongoClient, db, "messages", message)

}
