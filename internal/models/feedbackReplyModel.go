package models

type Reply struct {
	Content string `json:"feedback_reply"`
}

type EditReply struct {
	Content string `json:"content"`
	UserID  uint   `json:"userID"`
}

// struct to be stored in database
type NewFeedbackReply struct {
	Content    string
	FeedbackID string
	UserID     uint
}
