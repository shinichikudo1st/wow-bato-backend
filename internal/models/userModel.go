package models

type RegisterUser struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Barangay_ID string `json:"barangay"`
	Role        string `json:"role"`
	Contact     string `json:"contact"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStruct struct {
	ID          uint
	Password    string
	Role        string
	Barangay_ID uint
}

type UserProfile struct {
	ID        uint
	Email     string
	FirstName string
	LastName  string
	Role      string
	Contact   string
}
