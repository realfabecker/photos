package container

import (
	"context"
	"fmt"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/realfabecker/photos/internal/common/cache"
	"github.com/realfabecker/photos/internal/common/dotenv"
	"github.com/realfabecker/photos/internal/common/jwt"
	"log"

	payrep "github.com/realfabecker/photos/internal/adapters/photos/repositories"
	phtsrv "github.com/realfabecker/photos/internal/adapters/photos/services"
	usrsrv "github.com/realfabecker/photos/internal/adapters/users/services"
	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
	corsrv "github.com/realfabecker/photos/internal/core/services"
	"github.com/realfabecker/photos/internal/handlers/http"
	"github.com/realfabecker/photos/internal/handlers/http/routes"
	"go.uber.org/dig"
)

var Container dig.Container

func init() {
	Container = *dig.New()
	if err := reg3(); err != nil {
		log.Fatalf("unable to register services %v", err)
	}
}

func reg3() error {
	if err := Container.Provide(func() (*cordom.Config, error) {
		cnf := &cordom.Config{}
		if err := dotenv.Unmarshal(cnf); err != nil {
			return nil, err
		}
		return cnf, nil
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(cnf *cordom.Config) (*dynamodb.Client, error) {
		env, err := awsconf.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		return dynamodb.NewFromConfig(env), nil
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(cnf *cordom.Config) (*s3.Client, error) {
		env, err := awsconf.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		return s3.NewFromConfig(env), nil
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(cnf *cordom.Config, client *s3.Client) (corpts.MidiaBucket, error) {
		if cnf.BucketName == "" {
			return nil, fmt.Errorf("bucket name is required for midia bucket")
		}
		return phtsrv.NewS3MidiaSigner(cnf.BucketName, "photos", client), nil
	}); err != nil {
		return err
	}

	if err := Container.Provide(func() corpts.CacheHandler {
		return cache.NewFileCache()
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(cache corpts.CacheHandler) corpts.JwtHandler {
		return jwt.NewJwtHandler(cache)
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(cnf *cordom.Config) (*cognito.Client, error) {
		env, err := awsconf.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		return cognito.NewFromConfig(env), nil
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(d *dynamodb.Client, cnf *cordom.Config) (corpts.PhotoRepository, error) {
		return payrep.NewWalletDynamoDBRepository(d, cnf.DynamoDBTableName, cnf.AppName)
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(r corpts.PhotoRepository, m corpts.MidiaBucket) corpts.PhotoService {
		return corsrv.NewPhotoService(r, m)
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(
		walletConfig *cordom.Config,
		jwtHandler corpts.JwtHandler,
	) corpts.AuthService {
		return usrsrv.NewCognitoAuthService(walletConfig.CognitoJwkUrl, jwtHandler)
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(
		r corpts.PhotoRepository,
		s corpts.PhotoService,
		t corpts.AuthService,
	) (*routes.PhotoController, error) {
		return routes.NewPhotoController(r, s, t), nil
	}); err != nil {
		return err
	}

	if err := Container.Provide(func(
		appConfig *cordom.Config,
		photoController *routes.PhotoController,
		authService corpts.AuthService,
	) (corpts.HttpHandler, error) {
		return http.NewFiberHandler(
			appConfig,
			photoController,
			authService,
		), nil
	}); err != nil {
		return err
	}

	return nil
}
