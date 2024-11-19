package services

import (
	"errors"
	database "wow-bato-backend/internal"
	"wow-bato-backend/internal/models"
)

func RegisterUser(registerUser models.RegisterUser) error {
	db, err := database.ConnectDB()
	if err != nil {
		return err
	}

	hash, err := HashPassword(registerUser.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    	registerUser.Email,
		Password: 	hash,
		FirstName: 	registerUser.FirstName,
		LastName:  	registerUser.LastName,
		Role:      	registerUser.Role,
		Contact:   	registerUser.Contact,
	}

	result := db.Create(&user)
	return result.Error
}

func LoginUser(loginUser models.LoginUser) (models.User, error){
	db, err := database.ConnectDB()
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	result := db.Where("email = ?", loginUser.Email).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	if !CheckPassword(user.Password, loginUser.Password) {
		return models.User{}, errors.New("invalid password")
	}

	return user, nil
}
