package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"aws-secrets-example-go/providers"
)

// main runs the program from the CLI
func main() {
	const secretVar string = "SUPER_SECRET_VAR"
	const secretSsmKey string = "/app/some-secret"

	var sec providers.SecretProvider

	// Use environment variable secret if it's available, otherwise, use SSM
	// parameter store
	if _, ok := os.LookupEnv(secretVar); ok {
		sec = providers.NewEnvSecret(secretVar)
	} else {
		var err error
		sec, err = providers.NewSsmSecret(secretSsmKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		secretValue, err := sec.GetSecret()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(secretValue)
		time.Sleep(time.Second)
	}
}
