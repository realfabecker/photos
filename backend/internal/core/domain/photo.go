package domain

import (
	"fmt"
	"time"
)

// Photo @Description	Photo information
type Photo struct {
	PhotoId   string `json:"photoId" validate:"required" example:"2023050701GXEH91YBVV40C1FK50S1P0KC"`
	UserId    string `json:"userId" validate:"required" example:"e8ec3241-03b4-4aed-99d5-d72e1922d9b8"`
	FileName  string `json:"fileName" validate:"required,imagex_name" example:"image.jpg"`
	Title     string `json:"title" validate:"required" example:"Supermercado"`
	Url       string `json:"url" validate:"required,imagex_url" example:"https://images.com.br/image.jpg"`
	CreatedAt string `json:"createdAt" example:"2023-04-07T16:45:30Z"`
} // @name	Photo

// MidiaUpload @description Midia information
type MidiaUpload struct {
	Url string `json:"url" example:"https://images.com.br/image.jpg"`
} // @name MidiaUpload

type PhotoPagedDTOQuery struct {
	PagedDTOQuery
	CreatedAt *int32       `query:"created_at" validate:"omitempty" example:"2023"`
	Period    *PhotoPeriod `query:"period" validate:"omitempty,oneof=this_week this_month last_month next_month" example:"this_month"`
} //	@name	PhotoPagedDTOQuery

func (p PhotoPagedDTOQuery) GetPeriod() string {
	if p.Period != nil {
		return p.Period.Format()
	}
	if p.CreatedAt != nil {
		return fmt.Sprint(*p.CreatedAt)
	}
	return time.Now().Format("20060102")
}

type PhotoPeriod string

const (
	PhotoThisMonth PhotoPeriod = "this_month"
	PhotoLastMonth PhotoPeriod = "last_month"
)

func (p PhotoPeriod) String() string {
	switch p {
	case PhotoThisMonth:
		return "this_month"
	case PhotoLastMonth:
		return "last_month"
	}
	return "unknown"
}

func (p PhotoPeriod) Format() string {
	year, month, _ := time.Now().Date()
	switch p {
	case PhotoThisMonth:
		return time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Format("200601")
	case PhotoLastMonth:
		return time.Date(year, month-1, 1, 0, 0, 0, 0, time.Local).Format("200601")
	}
	return time.Now().Format("20060102")
}
