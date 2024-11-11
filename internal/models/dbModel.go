package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email   	string 	  `gorm:"unique;not null"`
	Password 	string 	  `gorm:"not null"`
	FirstName 	string 	  `gorm:"not null"`
	LastName 	string 	  `gorm:"not null"`
	Role     	string 	  `gorm:"not null"`
	Contact  	string 	  `gorm:"not null"`
	Barangay_ID uint 	  `gorm:"not null"`
	Barangay 	Barangay  `gorm:"foreignKey:Barangay_ID"`
	Feedbacks	[]Feedback `gorm:"foreignKey:UserID"`
}

type Barangay struct {
	gorm.Model
	Name     			string `gorm:"not null;unique"`
	City     			string `gorm:"not null"`
	Region   			string `gorm:"not null"`
	Users    			[]User `gorm:"foreignKey:Barangay_ID"`
	Projects 			[]Project `gorm:"foreignKey:Barangay_ID"`
	Budget_Categories 	[]Budget_Category `gorm:"foreignKey:Barangay_ID"`

}

type Budget_Category struct {
	gorm.Model
	Name     			string `gorm:"not null"`
	Description 		string `gorm:"type:text"`
	Barangay_ID 			uint   `gorm:"not null"`
	Barangay 			Barangay `gorm:"foreignKey:Barangay_ID"`
	Budget_Items 		[]Budget_Item `gorm:"foreignKey:CategoryID"`
	Projects 			[]Project `gorm:"foreignKey:CategoryID"`
}

type Budget_Item struct {
	gorm.Model
	Name     			string `gorm:"not null"`
	Amount_Allocated 	float64 `gorm:"not null"`
	Amount_Spent 		float64 `gorm:"not null"`
	Description 		string `gorm:"type:text"`
	Status 				string `gorm:"not null"` //pending, approved, rejected
	Approval_Date 		*time.Time //Nullable, set when approved
	CategoryID 			uint `gorm:"not null"`
	Category 			Budget_Category `gorm:"foreignKey:CategoryID"`
	ProjectID 			uint `gorm:"not null"`
	Project 			Project `gorm:"foreignKey:ProjectID"`

}

type Project struct {
	gorm.Model
	Name string `gorm:"not null"`
	Description string `gorm:"type:text"`
	StartDate time.Time `gorm:"not null"`
	EndDate *time.Time `gorm:"not null"`
	Status string `gorm:"not null"` //planned, ongoing, completed
	Barangay_ID uint `gorm:"not null"`
	Barangay Barangay `gorm:"foreignKey:Barangay_ID"`
	CategoryID uint `gorm:"not null"`
	Category Budget_Category `gorm:"foreignKey:CategoryID"`
	Budget_Items []Budget_Item `gorm:"foreignKey:ProjectID"`
}

type Feedback struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	Rating int `gorm:"not null"`
	UserID uint `gorm:"not null"`
	User User `gorm:"foreignKey:UserID"`
}

type Audit_Log struct {
	gorm.Model
}
