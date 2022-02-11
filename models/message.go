package models

type Message struct {
	ID             int    `json:"id,omitempty"`
	SenderUsername string `json:"sender_username"`
	TargetUsername string `json:"target_username"`
	Body           string `json:"body"`
}

