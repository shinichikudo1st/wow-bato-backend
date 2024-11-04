package models

type RegisterUser struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Role        string `json:"role"`
	Contact     string `json:"contact"`
	Barangay_ID uint   `json:"barangay_ID"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddBarangay struct {
	Name   string `json:"name"`
	City   string `json:"city"`
	Region string `json:"region"`
}

type DeleteBarangay struct {
	Barangay_ID uint `json:"barangay_ID"`
}

type UpdateBarangay struct {
	Barangay_ID uint `json:"barangay_ID"`
}
