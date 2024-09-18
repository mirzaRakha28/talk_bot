package request

// SendMessageToBotSubscriber represents the structure for sending a message to a bot subscriber
type SendMessageToBotSubscriberRequest struct {
	EmployeeCode string        `json:"employee_code"`
	Message      MessageSingle `json:"message"`
}

// SendMessageToBotGroupRequest represents the structure for sending a message to a bot group
type SendMessageToBotGroupRequest struct {
	GroupID string       `json:"group_id"`
	Message MessageGroup `json:"message"`
}

// Message represents the message details
type MessageGroup struct {
	Tag             string    `json:"tag"`
	Text            TextGroup `json:"text"`
	QuotedMessageID string    `json:"quoted_message_id,omitempty"`
	ThreadID        string    `json:"thread_id,omitempty"`
}

// Text represents the text content of the message
type TextGroup struct {
	Format  int    `json:"format"`
	Content string `json:"content"`
}

// Message represents the message details
type MessageSingle struct {
	Tag  string     `json:"tag"`
	Text TextSingle `json:"text"`
}

// Text represents the text content of the message
type TextSingle struct {
	Format  int    `json:"format"`
	Content string `json:"content"`
}
