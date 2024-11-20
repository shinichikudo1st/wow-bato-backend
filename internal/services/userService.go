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

func LoginUser(loginUser models.LoginUser) (models.UserStruct, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return models.UserStruct{}, err
	}

	var user models.UserStruct
	if err := db.Model(&models.User{}).
		Select("id, password, role").
		Where("email = ?", loginUser.Email).
		Scan(&user).Error; err != nil {
		return models.UserStruct{}, errors.New("invalid email or password")
	}

	if !CheckPassword(user.Password, loginUser.Password) {
		return models.UserStruct{}, errors.New("invalid email or password")
	}

	return user, nil
}

func GetUserProfile(userID uint) (models.UserProfile, error) {
	db, err := database.ConnectDB()
	if err != nil {
		return models.UserProfile{}, err
	}

	var userProfile models.UserProfile
	if err := db.Model(&models.User{}).
		Select("id, email, first_name, last_name, role, contact").
		Where("id = ?", userID).
		Scan(&userProfile).Error; err != nil {
		return models.UserProfile{}, err
	}

	return userProfile, nil
}
