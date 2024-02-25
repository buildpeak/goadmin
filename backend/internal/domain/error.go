package domain

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
}

func NewUserNotFoundError(id string) UserNotFoundError {
	return UserNotFoundError{
		notFoundError{
			baseError{
				Type:    "user_not_found",
				Message: "user not found",
			},
			"User",
		},
	}
}
