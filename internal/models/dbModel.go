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
	BarangayID 	uint 	  `gorm:"not null"`
	Barangay 	Barangay  `gorm:"foreignKey:BarangayID"`
	
}

type Barangay struct {
	gorm.Model
	Name     			string `gorm:"not null;unique"`
	City     			string `gorm:"not null"`
	Region   			string `gorm:"not null"`
	Users    			[]User `gorm:"foreignKey:BarangayID"`
	Projects 			[]Project `gorm:"foreignKey:BarangayID"`
	Budget_Categories 	[]Budget_Category `gorm:"foreignKey:BarangayID"`

}

type Budget_Category struct {
	gorm.Model
	Name     			string `gorm:"not null"`
	Description 		string `gorm:"type:text"`
	BarangayID 			uint   `gorm:"not null"`
	Barangay 			Barangay `gorm:"foreignKey:BarangayID"`
	Budget_Items 		[]Budget_Item `gorm:"foreignKey:Budget_CategoryID"`

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
	Projects 			[]Project `gorm:"many2many:project_budget_items"`

}

type Join_Project_Budget_Item struct {
	gorm.Model
	ProjectID uint `gorm:"primaryKey"`
	BudgetItemID uint `gorm:"primaryKey"`
	Project Project `gorm:"foreignKey:ProjectID"`
	BudgetItem Budget_Item `gorm:"foreignKey:BudgetItemID"`
}

type Project struct {
	gorm.Model
	Name string `gorm:"not null"`
	Description string `gorm:"type:text"`
	StartDate time.Time `gorm:"not null"`
	EndDate time.Time `gorm:"not null"`
	Status string `gorm:"not null"` //planned, ongoing, completed
	BarangayID uint `gorm:"not null"`
	Barangay Barangay `gorm:"foreignKey:BarangayID"`
	Budget_Items []Budget_Item `gorm:"many2many:project_budget_items"`
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
