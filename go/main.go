package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type secret struct {
	key     string
	value   string
	timeout int
	expiry  int
}

func newSecret(key string) secret {
	sec := new(secret)
	sec.key = key
	// 10 second cache of secret. Should default to a higher value and be configurable in production
	sec.timeout = 10
	return *sec
}

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

func (s *secret) isSecretValid() bool {
	return int(time.Now().Unix()) <= s.expiry
}

func (s *secret) resetExpiry() {
	s.expiry = int(time.Now().Unix()) + 10
}

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

func main() {
	sec := newSecret("/app/some-secret")
	for {
		sec.getSecret()
		fmt.Println(sec.value)
		time.Sleep(time.Second)
	}
}
