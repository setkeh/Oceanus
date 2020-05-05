package main

type image struct {
	ID  string `json:"id"`
	Src string `json:"src"`
	URL string `json:"url"`
	B64 string `json:"b64"`
}

type photo struct {
	ID  string `json:"id"`
	Src string `json:"src"`
	URL string `json:"url"`
}
