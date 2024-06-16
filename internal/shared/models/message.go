package models

type Message struct {
	// timestamp of the message
	Timestamp int64 `json:"timestamp"`

	// message content
	// TODO: expand the content to include more fields
	Content string `json:"content"`
}
