package db
import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"fmt"
)

var DB *gorm.DB

func migrations(){
	DB.AutoMigrate(&User{})
}

func New(){
	db , err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	DB = db

	if err != nil {
		panic("failed to connect database")
	} 

	migrations()
	
	fmt.Println("✅ - DB creation")
	fmt.Println("✅ - migrations are done ! ")
}


