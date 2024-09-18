package tokernservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"seatalk-bot/internal/config"
	"seatalk-bot/internal/constants"
	"seatalk-bot/models/response"
	"time"
)

// TokenService interacts with the Seatalk API for token management
type TokenService struct {
	config          *config.Config
	accessToken     string
	tokenExpireTime time.Time
}

// NewTokenService creates a new TokenService
func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{config: cfg}
}

// GetToken retrieves a new access token from the Seatalk API
func (s *TokenService) GetToken() (string, error) {
	url := s.config.AuthURL // Use Auth URL from config
	payload := map[string]string{
		"app_id":     s.config.AppID,
		"app_secret": s.config.AppSecret,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", errors.New(constants.ErrFailedToMarshalPayload)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", errors.New(constants.ErrFailedToCreateRequest)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New(constants.ErrFailedToExecuteRequest)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(constants.ErrFailedToGetToken + ": " + resp.Status)
	}

	var tokenResponse response.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", errors.New(constants.ErrFailedToDecodeResponse)
	}

	// Check if the response code is successful
	if tokenResponse.Code != 0 {
		return "", errors.New(constants.ErrApiError)
	}

	// Store the access token and its expiration time
	s.accessToken = tokenResponse.AppAccessToken
	s.tokenExpireTime = time.Unix(tokenResponse.Expire, 0)

	return s.accessToken, nil
}

// RefreshToken refreshes the access token if it's expired
func (s *TokenService) RefreshToken() (string, error) {
	// Check if the token is about to expire (e.g., 5 minutes before expiration)
	if time.Now().Add(5 * time.Minute).Before(s.tokenExpireTime) {
		return s.accessToken, nil // No need to refresh
	}

	// Call GetToken to refresh the token
	return s.GetToken()
}
