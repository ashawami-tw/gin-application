package utils

type Response struct {
	Message string
	Status  int
}

type Error struct {
	Error    string
	Response Response
}

func LogError(statusCode int, error, message string) *Error {
	return &Error{
		Error: error,
		Response: Response{
			Message: message,
			Status:  statusCode,
		},
	}
}
func ValidationError(tag string) string {
	switch tag {
	case "required":
		return "This is required field"
	case "email":
		return "Please provide valid email"
	}
	return ""
}
