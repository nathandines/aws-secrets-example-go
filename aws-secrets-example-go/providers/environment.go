package providers

import (
	"fmt"
	"os"
)

// EnvSecret is a secret which can be obtained from an environment variable
type EnvSecret struct {
	key   string // the name of the environment variable
	value string // the secret
}

// NewEnvSecret will initialise and return a secret from the environment
func NewEnvSecret(key string) *EnvSecret {
	sec := new(EnvSecret)
	sec.key = key
	return sec
}

// GetSecret returns the value of the environment variable
func (s *EnvSecret) GetSecret() (string, error) {
	if result, ok := os.LookupEnv(s.key); ok {
		return result, nil
	}
	return "", fmt.Errorf("Environment variable '%s' is not available", s.key)
}
