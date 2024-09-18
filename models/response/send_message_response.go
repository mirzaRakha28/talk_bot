package response

type SendMessageToBotSubscriberResponse struct {
	Code int `json:"code"`
}

type SendMessageToBotGroupResponse struct {
	Code      int    `json:"code"`
	MessegeId string `json:"message_id"`
}
