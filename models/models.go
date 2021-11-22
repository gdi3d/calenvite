package models

// Settings ...
type Settings struct {
	Mailgun       Mailgun
	SMTP          SMTP
	SenderAddress string
	SendUsing     string
	ServicePort   string
}

// Mailgun ...
type Mailgun struct {
	Domain    string
	SecretKey string
}

// SMTP ...
type SMTP struct {
	Host     string
	Port     string
	User     string
	Password string
}

// APIError ...
type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// ErrorFields ...
type ErrorFields struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// APIResponse ...
type APIResponse struct {
	Message     string        `json:"message"`
	StatusCode  int           `json:"status_code"`
	ErrorFields []ErrorFields `json:"error_fields"`
}

// RequestPayload ...
type RequestPayload struct {
	Users        []*InviteUser  `json:"users" validate:"required,dive"`
	Invitation   *CalendarEvent `json:"invitation" validate:"omitempty"`
	EmailSubject string         `json:"email_subject" validate:"required"`
	EmailBody    string         `json:"email_body" validate:"required"`
	EmailIsHTML  bool           `json:"email_is_html"`
}

// InviteUser ..
type InviteUser struct {
	FullName string `json:"full_name"`
	Email    string `json:"email" validate:"required,email"`
}

// CalendarEvent ...
type CalendarEvent struct {
	StartAt           string `json:"start_at" validate:"required"`
	EndAt             string `json:"end_at" validate:"required"`
	EventSummary      string `json:"summary"`
	Description       string `json:"description"`
	Location          string `json:"location"`
	OrganizerFullName string `json:"organizer_full_name" validate:"required"`
	OrganizerEmail    string `json:"organizer_email" validate:"required,email"`
}
