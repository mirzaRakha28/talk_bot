package request

type EventCallbackRequest struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Timestamp int64  `json:"timestamp"`
	AppID     string `json:"app_id"`
	Event     struct {
		SeaTalkChallenge string `json:"seatalk_challenge"`
		EmployeeCode     string `json:"employee_code"`
		GroupID          string `json:"group_id"`
		Message          struct {
			MessageID       string `json:"message_id"`
			QuotedMessageID string `json:"quoted_message_id"`
			ThreadID        string `json:"thread_id"`
			Sender          struct {
				SeatalkID    string `json:"seatalk_id"`
				EmployeeCode string `json:"employee_code"`
			} `json:"sender"`
			MessageSentTime int64  `json:"message_sent_time"`
			Tag             string `json:"tag"`
			Text            struct {
				Content       string `json:"content"`
				PlainText     string `json:"plain_text"`
				MentionedList []struct {
					Username  string `json:"username"`
					SeatalkID string `json:"seatalk_id"`
				} `json:"mentioned_list"`
			} `json:"text"`
		} `json:"message"`
		Group struct {
			GroupID       string `json:"group_id"`
			GroupName     string `json:"group_name"`
			GroupSettings struct {
				ChatHistoryForNewMembers string `json:"chat_history_for_new_members"`
				CanNotifyWithAtAll       bool   `json:"can_notify_with_at_all"`
				CanViewMemberList        bool   `json:"can_view_member_list"`
			} `json:"group_settings"`
		} `json:"group"`
		Inviter struct {
			SeaTalkID    string `json:"seatalk_id"`
			EmployeeCode string `json:"employee_code"`
		} `json:"inviter"`
	} `json:"event"`
}
type EventCallback struct {
	SeatalkChallenge string `json:"seatalk_challenge"`
}
