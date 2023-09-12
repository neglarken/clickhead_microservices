package model

type Item struct {
	Id    int    `json:"id" db:"id"`
	Info  string `json:"info" db:"info"`
	Price int    `json:"price" db:"price"`
}
