package models

type AddBarangay struct {
	Name   string `json:"name"`
	City   string `json:"city"`
	Region string `json:"region"`
}

type DeleteBarangay struct {
	Barangay_ID uint `json:"barangay_ID"`
}

type UpdateBarangay struct {
	Barangay_ID uint   `json:"barangay_ID"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Region      string `json:"region"`
}