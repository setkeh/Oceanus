package models

import (
	"image"
)

type Image struct {
	ID    string      `json:"id"`
	Src   string      `json:"src"`
	URL   string      `json:"url"`
	Image image.Image `json:"image"`
}

type Photo struct {
	ID  string `json:"id"`
	Src string `json:"src"`
	URL string `json:"url"`
}
