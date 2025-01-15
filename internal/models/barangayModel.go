package models

type AddBarangay struct {
	Name   string `json:"name"`
	City   string `json:"city"`
	Region string `json:"region"`
}

// CHANGE LATER TO NOT PUT BRGY. ID IN THE JSON BODY
// BUT AS A PARAMETER IN THE URL
type DeleteBarangay struct {
	Barangay_ID uint `json:"barangay_ID"`
}

type UpdateBarangay struct {
	Name   string `json:"name"`
	City   string `json:"city"`
	Region string `json:"region"`
}

// used for displaying brgy. information
type AllBarangayResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	City   string `json:"city"`
	Region string `json:"region"`
}

// returned to the client for selection inputs
type OptionBarangay struct {
    ID      uint `json:"id"`
    Name    string `json:"name"`
}
