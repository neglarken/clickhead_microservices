package model

type User struct {
	Id                  int    `json:"id" db:"id"`
	Login               string `json:"login" db:"login"`
	UnencryptedPassword string `json:"password"`
	Password            string `db:"hashed_password" json:"-"`
}
