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

func LoginUser(loginUser models.LoginUser) (models.UserResponse, error){
	db, err := database.ConnectDB()
	if err != nil {
		return models.UserResponse{}, err
	}

	var user models.UserStruct
	result := db.Model(&models.User{}).
        Select("id, email, password, first_name, last_name, role, contact").
        Where("email = ?", loginUser.Email).
        First(&user)

    if result.Error != nil {
        return models.UserResponse{}, errors.New("invalid email or password")
    }

	if !CheckPassword(user.Password, loginUser.Password) {
		return models.UserResponse{}, errors.New("invalid password")
	}

	return models.UserResponse{
		ID: user.ID,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Role: user.Role,
		Contact: user.Contact,
	}, nil
}
