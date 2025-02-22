// This is a simple docker secrets plugin that uses the 1Password Connect API to fetch secrets.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/docker/go-plugins-helpers/secrets"
)

// OPSecretsDriver implements the secrets.Driver interface
type OPSecretsDriver struct {
	op        connect.Client
	vaultName string
}

// Get gets a secret from a remote secret store
func (d OPSecretsDriver) Get(req secrets.Request) secrets.Response {
	log.Printf("Secret request: %#v", req)
	var itemName string
	var ok bool
	if itemName, ok = req.SecretLabels["item"]; !ok {
		err := fmt.Errorf("label 'item' missing in secret")
		log.Println(err)
		return secrets.Response{Err: err.Error(), DoNotReuse: true}
	}

	item, err := d.op.GetItem(itemName, d.vaultName)
	if err != nil {
		e := fmt.Errorf("error fetching item: %v", err)
		log.Println(e)
		return secrets.Response{Err: e.Error(), DoNotReuse: true}
	}

	value := item.GetValue(req.SecretName)
	if value == "" {
		err := fmt.Errorf("value not found in item")
		log.Println(err)
		return secrets.Response{Err: err.Error(), DoNotReuse: true}
	}

	return secrets.Response{Value: []byte(value), DoNotReuse: true}
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	vaultName := os.Getenv("OP_VAULT_NAME")
	if vaultName == "" {
		log.Fatalln("OP_VAULT_NAME environment variable is not set")
	}

	// Create a OPConnect client
	onePasswordClient, err := connect.NewClientFromEnvironment()
	if err != nil {
		log.Fatalf("Error creating one password client: %v", err)
	}

	// Create a new driver
	driver := OPSecretsDriver{onePasswordClient, vaultName}

	handler := secrets.NewHandler(driver)

	if err := handler.ServeUnix("/run/docker/plugins/opsd.sock", 0); err != nil {
		log.Fatalf("Error serving plugin: %v", err)
	}
}
