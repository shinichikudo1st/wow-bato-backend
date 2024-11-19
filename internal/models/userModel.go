package models

type RegisterUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	Contact   string `json:"contact"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStruct struct {
	ID        uint
	Email     string
	Password  string
	FirstName string
	LastName  string
	Role      string
	Contact   string
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	Contact   string `json:"contact"`
}
