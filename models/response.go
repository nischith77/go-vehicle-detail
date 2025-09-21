package models

type Collection struct {
	Url   string `json:"url"`
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Total int    `json:"total"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
	First string `json:"first"`
	Last  string `json:"last"`
}
type Data struct {
	Id     int    `json:"id"`
	MakeId int    `json:"make_id"`
	Make   string `json:"make"`
	Name   string `json:"name"`
}

type Response struct {
	Collection Collection `json:"collection"`
	Data       []Data     `json:"data"`
}
