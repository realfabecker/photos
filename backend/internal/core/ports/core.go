package ports

import (
	"github.com/golang-jwt/jwt/v5"
)

type HttpHandler interface {
	Register() error
	Listen(port string) error
	GetApp() interface{}
}

type JwtHandler interface {
	Decode(token string) (*jwt.RegisteredClaims, error)
}

type AuthService interface {
	Verify(token string) (*jwt.RegisteredClaims, error)
}
