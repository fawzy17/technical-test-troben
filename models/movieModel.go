package models

type MovieDetailModel struct {
	Name     string `json:"Title"`
	Year     string `json:"Year"`
	Poster   string `json:"Poster"`
	Genre    string `json:"Genre"`
	Type     string `json:"Type"`
	Director string `json:"Director"`
}
type MovieModel struct {
	Name     string `json:"Title"`
	Year     string `json:"Year"`
	Poster   string `json:"Poster"`
	Type     string `json:"Type"`
}
