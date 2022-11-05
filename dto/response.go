package dto

type Response struct {
	Message string `json:"message"`
	Errors  any    `json:"errors"`
	Data    any    `json:"data"`
}

func BuildResponse(message string, data any) Response {
	return Response{
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

func BuildErrorResponse(message string, err error) Response {
	return Response{
		Message: message,
		Errors:  err.Error(),
		Data:    nil,
	}
}

type Login struct {
	Username string `json:"username"`
	Nama string `json:"nama"`
	Token string `json:"token"`
}