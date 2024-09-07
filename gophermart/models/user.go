package models


type User struct {
	ID uint64
	Login string
	HashedPassword string `db:"hashed_password"`
	Balance Money
}
