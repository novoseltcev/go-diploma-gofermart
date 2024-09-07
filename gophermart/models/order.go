package models

import "time"

type Order struct {
	Number		string
	Status		string
	Accrual		*Money
	UserID		uint64
	UploadedAt 	time.Time
}
