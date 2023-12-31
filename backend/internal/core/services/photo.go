package services

import (
	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
)

type PhotoService struct {
	PhotoRepository corpts.PhotoRepository
	MidiaBucket     corpts.MidiaBucket
}

func NewPhotoService(r corpts.PhotoRepository, m corpts.MidiaBucket) corpts.PhotoService {
	return &PhotoService{PhotoRepository: r, MidiaBucket: m}
}

func (s *PhotoService) ListPhotos(email string, q cordom.PhotoPagedDTOQuery) (*cordom.PagedDTO[cordom.Photo], error) {
	p, err := s.PhotoRepository.ListPhotos(email, q)
	if err != nil {
		return nil, err
	}
	for i, v := range p.Items {
		if p.Items[i].Url, err = s.MidiaBucket.GetObjectUrl(v.Url, 1800); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (s *PhotoService) CreatePhoto(p *cordom.Photo) (*cordom.Photo, error) {
	p, err := s.PhotoRepository.CreatePhoto(p)
	if err != nil {
		return nil, err
	}
	if p.Url, err = s.MidiaBucket.GetObjectUrl(p.Url, 1800); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PhotoService) PutPhoto(p *cordom.Photo) (*cordom.Photo, error) {
	p, err := s.PhotoRepository.PutPhoto(p)
	if err != nil {
		return nil, err
	}
	if p.Url, err = s.MidiaBucket.GetObjectUrl(p.Url, 1800); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PhotoService) GetPhotoById(user string, photo string) (*cordom.Photo, error) {
	p, err := s.PhotoRepository.GetPhotoById(user, photo)
	if err != nil {
		return nil, err
	}
	if p.Url, err = s.MidiaBucket.GetObjectUrl(p.Url, 1800); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PhotoService) DeletePhoto(user string, photo string) error {
	p, err := s.PhotoRepository.GetPhotoById(user, photo)
	if err != nil {
		return err
	}
	if err := s.MidiaBucket.DeleteObject(p.Url); err != nil {
		return err
	}
	return s.PhotoRepository.DeletePhoto(user, photo)
}
