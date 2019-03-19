package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// secret stores secrets retrieved from SSM parameter store
type secret struct {
	key     string // SSM parameter store key in which the secret can be found
	value   string // the secret
	timeout int    // duration (in seconds) for a secret to be cached before it must be refreshed
	expiry  int    // UNIX time after which the current secret needs refreshing (now + timeout)
}

// newSecret will initialise and return a secret
func newSecret(key string) secret {
	sec := new(secret)
	sec.key = key
	// 10 second cache of secret. Should default to a higher value and be configurable in production
	sec.timeout = 10
	return *sec
}

// getSecret will return, and if necessary retrieve the secret defined in its SSM parameter store
// key
func (s *secret) getSecret() (string, error) {
	if s.isSecretValid() {
		return s.value, nil
	}
	log.Println("Refreshing secret...")
	if err := s.refreshSecret(); err != nil {
		return "", err
	}
	return s.value, nil
}

// isSecretValid will check if the current time is beyond the expiry of the secret. Being an int,
// expiry defaults to "0", so can be used immediately after the secret is created
func (s *secret) isSecretValid() bool {
	return int(time.Now().Unix()) <= s.expiry
}

// resetExpiry will save the expiry date, being the current time, plus the timeout defined when the
// secret was initialised
func (s *secret) resetExpiry() {
	s.expiry = int(time.Now().Unix()) + s.timeout
}

// refreshSecret will retrieve a secret from SSM parameter store, and reset the expiry for the
// secret
func (s *secret) refreshSecret() error {
	sess := session.Must(
		session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable}))

	ssmClient := ssm.New(sess)

	response, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(s.key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}
	s.resetExpiry()
	s.value = *response.Parameter.Value
	return nil
}

// main runs the program from the CLI
func main() {
	sec := newSecret("/app/some-secret")
	for {
		secretValue, err := sec.getSecret()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(secretValue)
		time.Sleep(time.Second)
	}
}
