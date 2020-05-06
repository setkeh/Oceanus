package models

type Image struct {
	ID  string `json:"id"`
	Src string `json:"src"`
	URL string `json:"url"`
	B64 string `json:"b64"`
}

type Photo struct {
	ID  string `json:"id"`
	Src string `json:"src"`
	URL string `json:"url"`
}
