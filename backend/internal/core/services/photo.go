package services

import (
	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
	"time"
)

type PhotoService struct {
	PhotoRepository corpts.PhotoRepository
	MidiaSigner     corpts.MidiaSigner
}

func NewPhotoService(r corpts.PhotoRepository, m corpts.MidiaSigner) corpts.PhotoService {
	return &PhotoService{PhotoRepository: r, MidiaSigner: m}
}

func (s *PhotoService) ListPhotos(email string, q cordom.PhotoPagedDTOQuery) (*cordom.PagedDTO[cordom.Photo], error) {
	return s.PhotoRepository.ListPhotos(email, q)
}

func (s *PhotoService) CreatePhoto(p *cordom.Photo) (*cordom.Photo, error) {
	return s.PhotoRepository.CreatePhoto(p)
}

func (s *PhotoService) PutPhoto(p *cordom.Photo) (*cordom.Photo, error) {
	return s.PhotoRepository.PutPhoto(p)
}

func (s *PhotoService) GetPhotoById(user string, photo string) (*cordom.Photo, error) {
	return s.PhotoRepository.GetPhotoById(user, photo)
}

func (s *PhotoService) DeletePhoto(user string, photo string) error {
	return s.PhotoRepository.DeletePhoto(user, photo)
}

func (s *PhotoService) GetUploadUrl(user string, name string) (string, error) {
	key := user + "/" + time.Now().Format("2006/01/02") + "/" + name
	return s.MidiaSigner.PutObjectUrl(key, 3600)
}
