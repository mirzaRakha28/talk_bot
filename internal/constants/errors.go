package constants

// Error messages
const (
	ErrFailedToMarshalPayload = "failed to marshal payload"
	ErrFailedToCreateRequest  = "failed to create HTTP request"
	ErrFailedToExecuteRequest = "failed to execute HTTP request"
	ErrFailedToGetToken       = "failed to get token"
	ErrFailedToDecodeResponse = "failed to decode response"
	ErrApiError               = "API returned an error"
	ErrorFileOpen             = "failed to open file"
	ErrorFileCreate           = "failed to create file"
	ErrorFileRead             = "error reading file"
	ErrorDateParse            = "failed to parse date"
	ErrorPICNotFound          = "PIC not found"
	ErrorPreviousPICNotFound  = "no previous PIC found"
	ErrorInvalidDateFormat    = "invalid date format"
	ErrorWriteSchedule        = "failed to write schedule"
)
