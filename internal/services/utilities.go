package services

import (
	"errors"
	"net/http"
	"wow-bato-backend/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func BindJSON(c *gin.Context, obj interface{}) {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}
}

func SetSession(session sessions.Session, user models.UserStruct){
	session.Set("user_id", user.ID)
	session.Set("user_role", user.Role)
	session.Set("barangay_id", user.Barangay_ID)
	session.Set("barangay_name", user.Barangay_Name)
	session.Set("authenticated", true)
}

func CheckServiceError(c *gin.Context, err error){
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func CheckAuthentication(c *gin.Context) sessions.Session{
	session := sessions.Default(c)
	if session.Get("authenticated") != true {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return nil
	}

	return session
}