package providers

// SecretProvider is an interface which represents types which provide secrets
type SecretProvider interface {
	GetSecret() (string, error)
}
