package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// main runs the program from the CLI
func main() {
	const secretVar string = "SUPER_SECRET_VAR"
	const secretSsmKey string = "/app/some-secret"

	var sec secretProvider

	// Use environment variable secret if it's available, otherwise, use SSM
	// parameter store
	if _, ok := os.LookupEnv(secretVar); ok {
		sec = newEnvSecret(secretVar)
	} else {
		sec = newSsmSecret(secretSsmKey)
	}

	for {
		secretValue, err := sec.getSecret()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(secretValue)
		time.Sleep(time.Second)
	}
}
