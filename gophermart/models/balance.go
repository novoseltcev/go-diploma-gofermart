package models


type Balance struct {
	Current Money `json:"current"`
	Withdrawn Money `json:"withdrawn"`
}
