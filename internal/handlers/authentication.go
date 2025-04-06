package handlers

import (
	"net/http"
	"wow-bato-backend/internal/models"
	"wow-bato-backend/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	svc *services.UserService
}

func NewUserHandlers(svc *services.UserService) *UserHandlers {
	return &UserHandlers{svc: svc}
}

func (h *UserHandlers) RegisterUser(c *gin.Context) {
	var registerUser models.RegisterUser

	services.BindJSON(c, &registerUser)

	err := h.svc.RegisterUser(registerUser)

	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h *UserHandlers) LoginUser(c *gin.Context){
	var loginUser models.LoginUser

	services.BindJSON(c, &loginUser)

	user, err := h.svc.LoginUser(loginUser)
	services.CheckServiceError(c, err)

	session := sessions.Default(c)
	services.SetSession(session, user)

	err = session.Save()
	services.CheckServiceError(c, err)
	

    c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged in successfully", "sessionStatus": session.Get("authenticated"), "role": session.Get("user_role")})

}

func (h *UserHandlers) LogoutUser(c *gin.Context){

	session := sessions.Default(c)
	session.Clear()

	session.Options(sessions.Options{
		Path: "/",
		MaxAge: -1,
		HttpOnly: true,
	})

	err := session.Save()
	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

func (h *UserHandlers) CheckAuth(c *gin.Context){
	

	session := services.CheckAuthentication(c)

	c.IndentedJSON(http.StatusOK, gin.H{"sessionStatus": session.Get("authenticated"), "role": session.Get("user_role"), 
	"user_id": session.Get("user_id"), "barangay_id": session.Get("barangay_id"), "barangay_name": session.Get("barangay_name")})
}

func (h *UserHandlers) GetUserProfile(c *gin.Context){

	session := services.CheckAuthentication(c)

	userID := session.Get("user_id")

	userProfile, err := h.svc.GetUserProfile(userID.(uint))

	services.CheckServiceError(c, err)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User profile fetched successfully", "data": userProfile})
}
