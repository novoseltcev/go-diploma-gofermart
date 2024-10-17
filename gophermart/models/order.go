package models

import "time"

type Order struct {
	Number		string
	Status		string
	Accrual		*Money
	UserID		uint64		`db:"user_id"`
	UploadedAt 	time.Time	`db:"uploaded_at"`
}
