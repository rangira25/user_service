package tasks

// Keep these task type names stable across services.
const (
	// Sent by auth_service when a new user is created.
	TypeSendWelcomeEmail = "email:send_welcome"
)
