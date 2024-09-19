package eventcallback

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"seatalk-bot/internal/config"
	"seatalk-bot/internal/constants"
	"seatalk-bot/models/request"
	"seatalk-bot/models/response"

	"github.com/robfig/cron/v3"
)

// EventCallbackService handles event callbacks and scheduled tasks
type EventCallbackService struct {
	config *config.Config
	cron   *cron.Cron
}

// NewEventCallbackService creates a new EventCallbackService
func NewEventCallbackService(cfg *config.Config) *EventCallbackService {
	service := &EventCallbackService{
		config: cfg,
		cron:   cron.New(),
	}

	// Schedule jobs
	service.scheduleJobs()

	return service
}

// scheduleJobs schedules periodic tasks
func (s *EventCallbackService) scheduleJobs() {
	// Schedule a job to run at 12 AM every Tuesday in Jakarta time
	s.cron.AddFunc("25 14 * * 3", func() {
		s.performScheduledPIC()
	})
	s.cron.AddFunc("0 0 * * 5", func() {
		s.performScheduledReminderReturnRefund()
	})
	s.cron.Start()
}

// performScheduledReminderReturnRefund checks if it's 12 AM Friday in Jakarta and performs the task
// for reminding test return and refund
func (s *EventCallbackService) performScheduledReminderReturnRefund() {
	req := request.SendMessageToBotGroupRequest{
		GroupID: s.config.RegressionGroupID, // Replace with the actual group ID
		Message: request.MessageGroup{
			Tag: "Text",
			Text: request.TextGroup{
				Format:  1,
				Content: constants.ReminderReturnRefund, // Use the constructed data as the message content
			},
		},
	}

	// Send the message to the group
	if _, err := s.SendMessageToGroup(req); err != nil {
		log.Println("Failed to send message to group:", err)
	}
}

// performScheduledTask checks if it's 12 AM Tuesday in Jakarta and performs the task
func (s *EventCallbackService) performScheduledPIC() {
	filename := "stock_inventory_schedule.txt"
	// Read schedules from file
	schedules, err := ReadSchedules(filename)
	if err != nil {
		log.Println(err)
		return
	}
	// Get the start and end date of the current week
	startOfWeek, endOfWeek := GetCurrentWeekRange()

	// Display PICs for the current week
	data := "PICs for this week: \n"
	// Display PICs for the current week
	data += DisplayPICsWithinRange(schedules, startOfWeek, endOfWeek)
	schedules, err = ReadSchedules(filename)
	if err != nil {
		log.Println(err)
		return
	}
	// Display the full schedule
	data += DisplayFullSchedule(schedules)

	req := request.SendMessageToBotGroupRequest{
		GroupID: s.config.RegressionGroupID, // Replace with the actual group ID
		Message: request.MessageGroup{
			Tag: "Text",
			Text: request.TextGroup{
				Format:  1,
				Content: data, // Use the constructed data as the message content
			},
		},
	}

	// Send the message to the group
	if _, err := s.SendMessageToGroup(req); err != nil {
		log.Println("Failed to send message to group:", err)
	}
}

// HandleEventCallback is the HTTP handler for event callbacks
func (s *EventCallbackService) HandleEventCallback(w http.ResponseWriter, r *http.Request) {
	// Limit to POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming request
	var eventRequest request.EventCallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Prepare the response message based on the event type
	switch eventRequest.EventType {
	case "message_from_bot_subscriber":
		req := request.SendMessageToBotSubscriberRequest{
			EmployeeCode: eventRequest.Event.EmployeeCode,
			Message: request.MessageSingle{
				Tag: "Text",
				Text: request.TextSingle{
					Format:  1,
					Content: "Message received. How can I help?",
				},
			},
		}
		if _, err := s.SendMessageToSubscriber(req); err != nil {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}
	case "new_mentioned_message_from_group_chat":
		req := request.SendMessageToBotGroupRequest{
			GroupID: eventRequest.Event.GroupID,
			Message: request.MessageGroup{
				Tag: "Text",
				Text: request.TextGroup{
					Format:  1,
					Content: "Message received. How can I help?",
				},
			},
		}
		if _, err := s.SendMessageToGroup(req); err != nil {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsupported event type", http.StatusBadRequest)
		return
	}

	// Prepare the response
	response := response.EventCallbackResponse{
		SeatalkChallenge: eventRequest.Event.SeaTalkChallenge,
	}

	// Set the response header and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SendMessageToSubscriber sends a message to a subscriber using the Seatalk API
func (s *EventCallbackService) SendMessageToSubscriber(req request.SendMessageToBotSubscriberRequest) (response.SendMessageToBotSubscriberResponse, error) {
	apiURL := s.config.SingleChatUrl

	// Marshal the request into JSON
	requestBody, err := json.Marshal(req)
	if err != nil {
		return response.SendMessageToBotSubscriberResponse{}, errors.New(constants.ErrFailedToMarshalPayload)
	}

	// Create a new HTTP request
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return response.SendMessageToBotSubscriberResponse{}, errors.New(constants.ErrFailedToCreateRequest)
	}

	// Set the request headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return response.SendMessageToBotSubscriberResponse{}, errors.New(constants.ErrFailedToExecuteRequest)
	}
	defer resp.Body.Close()

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return response.SendMessageToBotSubscriberResponse{}, errors.New(constants.ErrFailedToDecodeResponse)
		}
		return response.SendMessageToBotSubscriberResponse{
			Code: errorResponse.Code,
		}, errors.New(constants.ErrApiError)
	}

	// Parse the response
	var response response.SendMessageToBotSubscriberResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, errors.New(constants.ErrFailedToDecodeResponse)
	}

	return response, nil
}

// SendMessageToGroup sends a message to a group
func (s *EventCallbackService) SendMessageToGroup(req request.SendMessageToBotGroupRequest) (response.SendMessageToBotGroupResponse, error) {
	apiURL := s.config.GroupChatUrl

	// Marshal the request into JSON
	requestBody, err := json.Marshal(req)
	if err != nil {
		return response.SendMessageToBotGroupResponse{}, errors.New(constants.ErrFailedToMarshalPayload)
	}

	// Create a new HTTP request
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return response.SendMessageToBotGroupResponse{}, errors.New(constants.ErrFailedToCreateRequest)
	}

	// Set the request headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return response.SendMessageToBotGroupResponse{}, errors.New(constants.ErrFailedToExecuteRequest)
	}
	defer resp.Body.Close()

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return response.SendMessageToBotGroupResponse{}, errors.New(constants.ErrFailedToDecodeResponse)
		}
		return response.SendMessageToBotGroupResponse{
			Code:      errorResponse.Code,
			MessegeId: "0",
		}, errors.New(constants.ErrApiError)
	}

	// Parse the response
	var groupResponse response.SendMessageToBotGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&groupResponse); err != nil {
		return response.SendMessageToBotGroupResponse{}, errors.New(constants.ErrFailedToDecodeResponse)
	}

	return groupResponse, nil
}
