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
	Name   string `json:"name"`
	City   string `json:"city"`
	Region string `json:"region"`
}

type AllBarangayResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	City   string `json:"city"`
	Region string `json:"region"`
}

type OptionBarangay struct {
    ID      uint `json:"id"`
    Name    string `json:"name"`
}
