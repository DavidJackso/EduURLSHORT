package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func OK() Response {
	return Response{StatusOk, ""}
}

func Error(msg string) Response {
	return Response{StatusError, msg}
}
