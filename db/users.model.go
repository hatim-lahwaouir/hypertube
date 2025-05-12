package db

import ( 
	"time"
)
// User model
type User struct {
	ID		uint64	`gorm:"primaryKey"` 
	Username  	string	`gorm:"size:40;unique;not null"`
	Email     	string	`gorm:"size:40;unique;not null"`
	FirstName 	string	`gorm:"size:40;not null"`
	LastName  	string	`gorm:"size:40;not null"`
	Password  	string	`gorm:"size:100;not null"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
}
