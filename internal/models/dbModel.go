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
	ProfilePicture  string    `gorm:"default:null"`
	Contact  	string 	  `gorm:"not null"`
	Barangay_ID *uint 	  `gorm:"default:null"`
	Barangay 	Barangay  `gorm:"foreignKey:Barangay_ID"`
	Feedbacks	[]Feedback `gorm:"foreignKey:UserID"`
	FeedbackReplies []FeedbackReply `gorm:"foreignKey:UserID"`
}

type Barangay struct {
	gorm.Model
	Name     			string `gorm:"not null;unique"`
	City     			string `gorm:"not null"`
	Region   			string `gorm:"not null"`
    ImageURL            string `gorm:""`
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
	Projects 			[]Project `gorm:"foreignKey:CategoryID"`
}

type Project struct {
	gorm.Model
	Name string `gorm:"not null"`
	Description string `gorm:"type:text"`
	StartDate time.Time `gorm:"not null"`
	EndDate time.Time `gorm:"not null"`
	Status string `gorm:"not null"` //planned, ongoing, completed
	Barangay_ID uint `gorm:"not null"`
	Barangay Barangay `gorm:"foreignKey:Barangay_ID"`
	CategoryID uint `gorm:"not null"`
	Category Budget_Category `gorm:"foreignKey:CategoryID"`
    Feedbacks []Feedback `gorm:"foreignKey:ProjectID"`
	Budget_Items []Budget_Item `gorm:"foreignKey:ProjectID"`
}

type Budget_Item struct {
	gorm.Model
	Name     			string `gorm:"not null"`
	Amount_Allocated 	float64 `gorm:"not null"`
	Description 		string `gorm:"type:text"`
	Status 				string `gorm:"not null"` //pending, approved, rejected
	Approval_Date 		*time.Time //Nullable, set when approved
	ProjectID 			uint `gorm:"not null"`
	Project 			Project `gorm:"foreignKey:ProjectID"`
}

type Feedback struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	Role string `gorm:"not null"`
    ProjectID uint `gorm:"not null"`
    Project Project `gorm:"foreignKey:ProjectID"`
	UserID uint `gorm:"not null"`
	User User `gorm:"foreignKey:UserID"`
	FeedbackReplies []FeedbackReply `gorm:"foreignKey:FeedbackID"`
}

type FeedbackReply struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	FeedbackID uint `gorm:"not null"`
	Feedback Feedback `gorm:"foreignKey:FeedbackID"`
	UserID uint `gorm:"not null"`
	User User `gorm:"foreignKey:UserID"`
}
