package restfulapi

type RestfulAPIResponse struct {
	Status	int
	Data	interface{}
}

type SuccessMessage struct {
	Message		string
}

type ErrorMessage struct {
	ErrorCode	int8
	ErrorMsg	string
}

func Success(data interface{}) RestfulAPIResponse {
	return RestfulAPIResponse{0, data}
}

func SuccessMsg(msg string) RestfulAPIResponse {
	res := SuccessMessage {
		Message:	msg,
	}
	return RestfulAPIResponse{0, res}
}

func Error(code int8, msg string) RestfulAPIResponse {
	err := ErrorMessage {
		ErrorCode:	code,
		ErrorMsg:	msg,
	}
	return RestfulAPIResponse{1, err}
}