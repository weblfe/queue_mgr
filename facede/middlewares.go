package facede

type SecretStorage interface {
	GetSecretByAppID(appID string) string
}
