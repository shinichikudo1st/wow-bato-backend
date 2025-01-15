package models

// JSON struct for creating new user
type RegisterUser struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Barangay_ID string `json:"barangay"`
	Role        string `json:"role"`
	Contact     string `json:"contact"`
}

// JSON struct for logging in
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// struct used for storing session data
type UserStruct struct {
	ID            uint
	Password      string
	Role          string
	Barangay_ID   uint
	Barangay_Name string
}

// struct to be returned for user profile display
type UserProfile struct {
	ID        uint
	Email     string
	FirstName string
	LastName  string
	Role      string
	Contact   string
}
