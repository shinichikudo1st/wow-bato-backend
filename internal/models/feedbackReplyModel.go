package models

type Reply struct {
	Content string `json:"feedback_reply"`
}

// struct to be stored in database
type NewFeedbackReply struct {
	Content    string
	FeedbackID string
	UserID     uint
}
