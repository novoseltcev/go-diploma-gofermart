package models

import "time"

type Order struct {
	Number		string
	Status		string
	Accrual		Money
	UploadedAt 	time.Time
}
