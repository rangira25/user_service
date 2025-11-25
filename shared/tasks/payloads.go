package tasks

// WelcomeEmailPayload is the contract used by auth_service and notification worker.
// Keep JSON tags stable.
type WelcomeEmailPayload struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	ReqID  string `json:"req_id"`
}
