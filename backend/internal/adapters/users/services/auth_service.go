package users

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
)

type CognitoAuthService struct {
	cognitoClient   *cognitoidentityprovider.Client
	cognitoClientId string
	cognitoJwkUrl   string
	jwtHandler      corpts.JwtHandler
}

func NewCognitoAuthService(
	cognitoClientId string,
	cognitoJwkUrl string,
	cognitoClient *cognitoidentityprovider.Client,
	jwtHandler corpts.JwtHandler,
) corpts.AuthService {
	return &CognitoAuthService{cognitoClient, cognitoClientId, cognitoJwkUrl, jwtHandler}
}

func (u *CognitoAuthService) Login(email string, password string) (*cordom.UserToken, error) {
	auth, err := u.cognitoClient.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: aws.String(u.cognitoClientId),
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	})
	if err != nil {
		return nil, err
	}
	if auth.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		return &cordom.UserToken{
			AuthSession:   auth.Session,
			AuthChallenge: aws.String("NEW_PASSWORD_REQUIRED"),
		}, nil
	}
	if auth.AuthenticationResult.AccessToken == nil {
		return nil, errors.New("unexpected authorization error")
	}
	return &cordom.UserToken{
		AccessToken:  auth.AuthenticationResult.AccessToken,
		RefreshToken: auth.AuthenticationResult.RefreshToken,
	}, nil
}

func (u *CognitoAuthService) Change(email string, password string, session string) (*cordom.UserToken, error) {
	auth, err := u.cognitoClient.RespondToAuthChallenge(context.TODO(), &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: "NEW_PASSWORD_REQUIRED",
		ClientId:      &u.cognitoClientId,
		Session:       &session,
		ChallengeResponses: map[string]string{
			"USERNAME":     email,
			"NEW_PASSWORD": password,
		},
	})
	if err != nil {
		return nil, err
	}
	if auth.AuthenticationResult.AccessToken == nil {
		return nil, errors.New("unexpected authorization error")
	}
	return &cordom.UserToken{
		AccessToken:  auth.AuthenticationResult.AccessToken,
		RefreshToken: auth.AuthenticationResult.RefreshToken,
	}, nil
}

func (u *CognitoAuthService) Verify(token string) (*jwt.RegisteredClaims, error) {
	c, err := u.jwtHandler.VerifyWithKeyURL(token, u.cognitoJwkUrl)
	if err != nil {
		return nil, err
	}
	if time.Now().Unix() > c.ExpiresAt.Time.Unix() {
		return nil, errors.New("token expired")
	}
	return c, nil
}
