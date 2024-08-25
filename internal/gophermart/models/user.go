package models


type User struct {
	Id uint64
	Login string
	HashedPassword string `db:"hashed_password"`
	Balance Money
}
