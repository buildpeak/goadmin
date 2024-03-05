package domain

// baseError is a struct that represents the base error type for the API.
type baseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e baseError) Error() string {
	return e.Message
}

type notFoundError struct {
	baseError
	Resource string `json:"resource"`
}

type UserNotFoundError struct {
	notFoundError
	id string
}

func NewUserNotFoundError(usrID string) UserNotFoundError {
	return UserNotFoundError{
		notFoundError: notFoundError{
			baseError{
				Type:    "user_not_found",
				Message: "user not found",
			},
			"User",
		},
		id: usrID,
	}
}
