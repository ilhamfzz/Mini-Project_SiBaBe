package dto

type Response struct {
	Message string `json:"message"`
	Errors  any    `json:"errors"`
	Data    any    `json:"data"`
}

type EmptyObj struct{}

func BuildResponse(message string, data any) Response {
	return Response{
		Message: message,
		Errors:  EmptyObj{},
		Data:    data,
	}
}

func BuildErrorResponse(message string, err error) Response {
	return Response{
		Message: message,
		Errors:  err.Error(),
		Data:    EmptyObj{},
	}
}

type Login struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Token    string `json:"token"`
}
