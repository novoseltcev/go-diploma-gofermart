package models


type Balance struct {
	Balance Money `json:"balance"`
	Withdrawn Money `json:"withdrawn"`
}
