package main

type secretProvider interface {
	getSecret() (string, error)
}
