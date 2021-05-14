package token

import "time"

// Maker is an interface for manage token
type Maker interface {
	// CreateToken creates a new token from specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checkes if the token valid or not
	VerifyToken(token string) (*Payload, error)
}
