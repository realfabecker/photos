package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/realfabecker/photos/internal/adapters/common/validator"
	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
)

type PhotoController struct {
	repository corpts.PhotoRepository
	service    corpts.PhotoService
	auth       corpts.AuthService
}

func NewPhotoController(
	walletRepository corpts.PhotoRepository,
	walletService corpts.PhotoService,
	auth corpts.AuthService,
) *PhotoController {
	return &PhotoController{walletRepository, walletService, auth}
}

// ListPhotos list photos
//
//	@Summary		List photos
//	@Description	List photos
//	@Tags			Photos
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			limit		query		number	true	"Number of records"
//	@Param			page_token	query		string	false	"Pagination token"
//	@Param			due_date	query		string	true	"Photo due date"
//	@Success		200			{object}	cordom.ResponseDTO[cordom.PagedDTO[cordom.Photo]]
//	@Failure		400
//	@Failure		500
//	@Router			/photos [get]
func (w *PhotoController) ListPhotos(c *fiber.Ctx) error {
	q := cordom.PhotoPagedDTOQuery{}
	if err := c.QueryParser(&q); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	v := validator.NewValidator()
	if err := v.Struct(q); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, ok := c.Locals("user").(*jwt.RegisteredClaims)
	if !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code)
	}

	out, err := w.service.ListPhotos(user.Subject, q)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(cordom.ResponseDTO[cordom.PagedDTO[cordom.Photo]]{
		Status: "success",
		Data:   out,
	})
}

// GetPhotoById get a photo by its id
//
//	@Summary		Get photo by id
//	@Description	Get photo by id
//	@Tags			Photos
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		string	true	"Photo id"
//	@Success		200	{object}	cordom.ResponseDTO[cordom.Photo]
//	@Failure		400
//	@Failure		500
//	@Router			/photos/{photoId} [get]
func (w *PhotoController) GetPhotoById(c *fiber.Ctx) error {
	p := cordom.Photo{}
	if err := c.ParamsParser(&p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	v := validator.NewValidator()
	if err := v.StructPartial(p, "photoId"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, ok := c.Locals("user").(*jwt.RegisteredClaims)
	if !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code)
	}

	d, err := w.service.GetPhotoById(user.Subject, p.PhotoId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if d == nil && err == nil {
		return fiber.NewError(fiber.StatusNotFound)
	}

	return c.JSON(cordom.ResponseDTO[cordom.Photo]{
		Status: "success",
		Data:   d,
	})
}

// CreatePhoto create a new photo
//
//	@Summary		Create a photo
//	@Description	Create a new photo record
//	@Tags			Photos
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			request	body		cordom.Photo	true	"Photo payload"
//	@Success		200		{object}	cordom.ResponseDTO[cordom.Photo]
//	@Failure		400
//	@Failure		500
//	@Router			/photos [post]
func (w *PhotoController) CreatePhoto(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.RegisteredClaims)
	if !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code)
	}
	body := cordom.Photo{UserId: user.Subject}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	v := validator.NewValidator()
	if err := v.StructExcept(body, "PhotoId"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	p, err := w.service.CreatePhoto(&body)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(cordom.ResponseDTO[cordom.Photo]{
		Status: "success",
		Data:   p,
	})
}

// DeletePhoto delete a photo by its Id
//
//	@Summary		Delete photo
//	@Description	Delete photo
//	@Tags			Photos
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id	path		string	true	"Photo id"
//	@Success		200	{object}	cordom.EmptyResponseDTO
//	@Failure		400
//	@Failure		500
//	@Router			/photos/{photoId} [delete]
func (w *PhotoController) DeletePhoto(c *fiber.Ctx) error {
	p := cordom.Photo{}
	if err := c.ParamsParser(&p); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	v := validator.NewValidator()
	if err := v.StructPartial(p, "Id"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, ok := c.Locals("user").(*jwt.RegisteredClaims)
	if !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code)
	}

	photo, err := w.service.GetPhotoById(user.Subject, p.PhotoId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if photo == nil {
		return fiber.NewError(fiber.StatusNotFound, "Not Found")
	}

	if err := w.service.DeletePhoto(user.Subject, p.PhotoId); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(204)
}

// PutPhoto
//
//	@Summary		Put a photo
//	@Description	Update/Create a photo record
//	@Tags			Photos
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			id		path		string			true	"Photo id"
//	@Param			request	body		cordom.Photo	true	"Photo payload"
//	@Success		200		{object}	cordom.ResponseDTO[cordom.Photo]
//	@Failure		400
//	@Failure		500
//	@Router			/photos/{photoId} [put]
func (w *PhotoController) PutPhoto(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.RegisteredClaims)
	if !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code)
	}

	body := cordom.Photo{UserId: user.Subject}
	if err := c.ParamsParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	photo, err := w.service.GetPhotoById(user.Subject, body.PhotoId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if photo == nil {
		return fiber.NewError(fiber.StatusNotFound, "Not Found")
	}

	if err := c.BodyParser(&photo); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	v := validator.NewValidator()
	if err := v.Struct(photo); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = w.service.PutPhoto(photo)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(204)
}

// GetUploadUrl get a photo upload url
//
//	@Summary		Get photo upload url
//	@Description	Get photo upload url
//	@Tags			Bucket
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			file	query		string	true	"Filename"
//	@Success		200		{object}	cordom.ResponseDTO[cordom.MidiaUpload]
//	@Failure		400
//	@Failure		500
//	@Router			/bucket/upload-url [get]
func (w *PhotoController) GetUploadUrl(c *fiber.Ctx) error {
	q := struct {
		File string `json:"file" query:"file" validate:"required,imagex_name"`
	}{}

	if err := c.QueryParser(&q); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	v := validator.NewValidator()
	if err := v.StructPartial(q, "File"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, ok := c.Locals("user").(*jwt.RegisteredClaims)
	if !ok {
		return fiber.NewError(fiber.ErrUnauthorized.Code)
	}

	url, err := w.service.GetUploadUrl(user.Subject, q.File)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if url == "" && err == nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	return c.JSON(cordom.ResponseDTO[cordom.MidiaUpload]{
		Status: "success",
		Data: &cordom.MidiaUpload{
			Url: url,
		},
	})
}
