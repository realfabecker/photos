package domain

type UserToken struct {
	RefreshToken  *string `json:"RefreshToken,omitempty"`
	AccessToken   *string `json:"AccessToken,omitempty"`
	AuthChallenge *string `json:"AuthChallenge,omitempty"`
	AuthSession   *string `json:"AuthSession,omitempty"`
} // @name	UserToken

type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
} // @name	UserLoginDTO

type UserLoginChangeDTO struct {
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
	Session     string `json:"session" validate:"required"`
} // @name	UserLoginChangeDTO
