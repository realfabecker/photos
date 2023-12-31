package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/swagger"

	"github.com/realfabecker/photos/internal/handlers/http/docs"

	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
	"github.com/realfabecker/photos/internal/handlers/http/routes"
)

type Handler struct {
	app             *fiber.App
	photoConfig     *cordom.Config
	photoController *routes.PhotoController
	authService     corpts.AuthService
}

//	@title						Photos Rest API
//	@version					1.0
//	@description				Photos Rest API
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.htm
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Type 'Bearer ' and than your API token
func NewFiberHandler(
	photoConfig *cordom.Config,
	photoController *routes.PhotoController,
	authService corpts.AuthService,
) corpts.HttpHandler {

	// open api base project configuration (2)
	docs.SwaggerInfo.Host = photoConfig.AppHost
	docs.SwaggerInfo.Schemes = []string{"http"}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msgs := utils.StatusMessage(code)

			var ferr *fiber.Error
			if errors.As(err, &ferr) {
				code = ferr.Code
				msgs = ferr.Message
			}

			c.Status(code)
			return c.JSON(cordom.ResponseDTO[interface{}]{
				Status:  "error",
				Message: msgs,
				Code:    code,
			})
		},
	})
	return &Handler{
		app,
		photoConfig,
		photoController,
		authService,
	}
}

func (a *Handler) GetApp() interface{} {
	return a.app
}

func (a *Handler) Listen(port string) error {
	return a.app.Listen(":" + port)
}

func (a *Handler) Register() error {
	a.app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 30 * time.Second,
	}))

	a.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	a.app.Get("/docs/*", swagger.HandlerDefault)
	photos := a.app.Group("/photos")

	media := photos.Group("/media")
	media.Use(a.authHandler)
	media.Post("/", a.photoController.CreatePhoto)
	media.Get("/", a.photoController.ListPhotos)
	media.Get("/:photoId", a.photoController.GetPhotoById)
	media.Delete("/:photoId", a.photoController.DeletePhoto)
	media.Put("/:photoId", a.photoController.PutPhoto)
	return nil
}

func (a *Handler) authHandler(c *fiber.Ctx) error {
	auth := c.Get("authorization")
	if len(auth) < (len("bearer") + 1) {
		return fiber.NewError(fiber.ErrUnauthorized.Code)
	}
	u, err := a.authService.Verify(auth[len("bearer")+1:])
	if err != nil {
		return fiber.NewError(fiber.ErrUnauthorized.Code, err.Error())
	}
	c.Locals("user", u)
	return c.Next()
}
