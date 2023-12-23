package ports

import (
	"github.com/golang-jwt/jwt/v5"
	cordom "github.com/realfabecker/photos/internal/core/domain"
)

type AuthService interface {
	Login(email string, password string) (*cordom.UserToken, error)
	Change(email string, password string, session string) (*cordom.UserToken, error)
	Verify(token string) (*jwt.RegisteredClaims, error)
}

type UserRepository interface {
	GetUserByEmail(email string) (*cordom.User, error)
}

type UserService interface {
	GetUserByEmail(email string) (*cordom.User, error)
}
