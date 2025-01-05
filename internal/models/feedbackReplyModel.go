package models

type Reply struct {
	Content string `json:"feedback_reply"`
}

type NewFeedbackReply struct {
	Content    string
	FeedbackID string
	UserID     uint
}