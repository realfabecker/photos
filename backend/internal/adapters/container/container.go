package container

import (
	"context"
	"fmt"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/realfabecker/photos/internal/adapters/common/cache"
	"github.com/realfabecker/photos/internal/adapters/common/env"
	"github.com/realfabecker/photos/internal/adapters/common/jwt"

	payrep "github.com/realfabecker/photos/internal/adapters/photos/repositories"
	phtsrv "github.com/realfabecker/photos/internal/adapters/photos/services"
	usrrep "github.com/realfabecker/photos/internal/adapters/users/repositories"
	usrsrv "github.com/realfabecker/photos/internal/adapters/users/services"
	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
	corsrv "github.com/realfabecker/photos/internal/core/services"
	"github.com/realfabecker/photos/internal/handlers/http"
	"github.com/realfabecker/photos/internal/handlers/http/routes"
	"go.uber.org/dig"
)

var Container dig.Container

// init
func init() {
	Container = *dig.New()

	Container.Provide(func() (*cordom.Config, error) {
		cnf := &cordom.Config{}
		if err := env.Unmarshal(cnf); err != nil {
			return nil, err
		}
		return cnf, nil
	})

	Container.Provide(func(cnf *cordom.Config) (*dynamodb.Client, error) {
		env, err := awsconf.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		return dynamodb.NewFromConfig(env), nil
	})

	Container.Provide(func(cnf *cordom.Config) (*s3.Client, error) {
		env, err := awsconf.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		return s3.NewFromConfig(env), nil
	})

	Container.Provide(func(cnf *cordom.Config, client *s3.Client) (corpts.MidiaSigner, error) {
		if cnf.BucketName == "" {
			return nil, fmt.Errorf("bucket name is required for midia bucket")
		}
		return phtsrv.NewS3MidiaSigner(cnf.BucketName, "photos", client), nil
	})

	Container.Provide(func() corpts.CacheHandler {
		return cache.NewFileCache()
	})

	Container.Provide(func(cache corpts.CacheHandler) corpts.JwtHandler {
		return jwt.NewJwtHandler(cache)
	})

	Container.Provide(func(cnf *cordom.Config) (*cognito.Client, error) {
		env, err := awsconf.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}

		return cognito.NewFromConfig(env), nil
	})

	Container.Provide(func(d *dynamodb.Client, cnf *cordom.Config) (corpts.PhotoRepository, error) {
		return payrep.NewWalletDynamoDBRepository(d, cnf.DynamoDBTableName, cnf.AppName)
	})

	Container.Provide(func(r corpts.PhotoRepository, m corpts.MidiaSigner) corpts.PhotoService {
		return corsrv.NewPhotoService(r, m)
	})

	Container.Provide(func(
		walletConfig *cordom.Config,
		cognitoClient *cognito.Client,
		jwtHandler corpts.JwtHandler,
	) corpts.AuthService {
		return usrsrv.NewCognitoAuthService(
			walletConfig.CognitoClientId,
			walletConfig.CognitoJwkUrl,
			cognitoClient,
			jwtHandler,
		)
	})

	Container.Provide(func(
		r corpts.PhotoRepository,
		s corpts.PhotoService,
		t corpts.AuthService,
	) (*routes.PhotoController, error) {
		return routes.NewPhotoController(r, s, t), nil
	})

	Container.Provide(func(d *dynamodb.Client, cnf *cordom.Config) (corpts.UserRepository, error) {
		return usrrep.NewUserRepository(d, cnf.DynamoDBTableName, cnf.AppName)
	})

	Container.Provide(func(
		r corpts.UserRepository,
	) (corpts.UserService, error) {
		return corsrv.NewUserService(r), nil
	})

	Container.Provide(func(
		a corpts.AuthService,
		u corpts.UserService,
	) (*routes.AuthController, error) {
		return routes.NewAuthController(a, u), nil
	})

	Container.Provide(func(
		walletConfig *cordom.Config,
		walletController *routes.PhotoController,
		usersController *routes.AuthController,
		authService corpts.AuthService,
	) (corpts.HttpHandler, error) {
		return http.NewFiberHandler(
			walletConfig,
			walletController,
			usersController,
			authService,
		), nil
	})
}
