package database

import (
	"wow-bato-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=MniVAPBNHS123 dbname=wowBato port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Barangay{}, &models.Budget_Category{}, &models.Budget_Item{}, &models.Project{}, &models.Feedback{}, &models.FeedbackReply{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
