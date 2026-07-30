package entity

type BaseResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func (r BaseResponse) ToResponse(
	message string,
	status int,
	data, err interface{},
) BaseResponse {
	return BaseResponse{
		Message: message,
		Status:  status,
		Data:    data,
		Error:   err,
	}
}
