package credentials

var _ KeychainCredentialsProvider = (*Credentials)(nil)

// KeychainCredentialsProvider provides KeychainCredentials.
type KeychainCredentialsProvider interface {
	KeychainCredentials() KeychainCredentials
}

// KeychainCredentials provides KeychainCredentials.
func (c *Credentials) KeychainCredentials() KeychainCredentials {
	return c
}
