package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"user/internal/model"
)

func Init() *gorm.DB {
	dsn := "host=postgres user=postgres password=phong dbname=profile port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("AutoMigrate got err : %v\n", err)
	}

	return db
}
