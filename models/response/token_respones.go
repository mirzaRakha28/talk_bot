package response

// TokenResponse represents the response from the Seatalk API for token requests
type TokenResponse struct {
	Code           int    `json:"code"`
	AppAccessToken string `json:"app_access_token"`
	Expire         int64  `json:"expire"`
}
