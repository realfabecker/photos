package ports

import (
	cordom "github.com/realfabecker/photos/internal/core/domain"
)

type PhotoRepository interface {
	ListPhotos(email string, q cordom.PhotoPagedDTOQuery) (*cordom.PagedDTO[cordom.Photo], error)
	CreatePhoto(p *cordom.Photo) (*cordom.Photo, error)
	PutPhoto(p *cordom.Photo) (*cordom.Photo, error)
	GetPhotoById(user string, photo string) (*cordom.Photo, error)
	DeletePhoto(user string, photo string) error
}

type PhotoService interface {
	ListPhotos(email string, q cordom.PhotoPagedDTOQuery) (*cordom.PagedDTO[cordom.Photo], error)
	CreatePhoto(p *cordom.Photo) (*cordom.Photo, error)
	PutPhoto(p *cordom.Photo) (*cordom.Photo, error)
	GetPhotoById(user string, photo string) (*cordom.Photo, error)
	DeletePhoto(user string, photo string) error
	GetUploadUrl(user string, name string) (string, error)
}

type MidiaSigner interface {
	GetObjectUrl(name string, lifetime int64) (string, error)
	PutObjectUrl(name string, lifetime int64) (string, error)
}
