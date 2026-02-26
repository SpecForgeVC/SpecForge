package auth

type Config struct {
	SigningKey []byte
	Issuer     string
	Audience   string
	Algorithm  string // e.g., "HS256" or "RS256"
}
