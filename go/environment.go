package main

import (
	"fmt"
	"os"
)

type envSecret struct {
	key   string // the name of the environment variable
	value string // the secret
}

// newEnvSecret will initialise and return a secret from the environment
func newEnvSecret(key string) *envSecret {
	sec := new(envSecret)
	sec.key = key
	return sec
}

// getSecret returns the value of the environment variable
func (s *envSecret) getSecret() (string, error) {
	if result, ok := os.LookupEnv(s.key); ok {
		return result, nil
	}
	return "", fmt.Errorf("Environment variable '%s' is not available", s.key)
}
