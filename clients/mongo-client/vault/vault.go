package vault

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
)

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
		log.Fatal(err)
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

		color.Cyan(fmt.Sprintf("\nVault Request: %v [OK]  You have successfully authenticated with Vault and a MongoDB credential has been created and will last for %v seconds\n", resp.StatusCode, VaultJSONResponse["lease_duration"]))
		fmt.Println("")
		color.Cyan(fmt.Sprintf("\tUserName: %v\n\tPassword: %v", username, password))
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
