package helper

import "github.com/Arshia-Izadyar/Go-Ecommerce/src/api/validator"

type Response struct {
	Result     interface{}
	StatusCode int
	Success    bool
	Error      string
}

func GenerateResponse(res interface{}, code int, success bool) *Response {
	return &Response{
		Result:     res,
		StatusCode: code,
		Success:    success,
		Error:      "",
	}
}

func GenerateResponseWithError(err error, code int, success bool) *Response {
	return &Response{
		Result:     nil,
		StatusCode: code,
		Success:    success,
		Error:      err.Error(),
	}
}
func GenerateResponseWithValidationError(err error, code int, success bool) *Response {
	ve := validator.GetValidationError(err)
	return &Response{
		Result:     ve,
		StatusCode: code,
		Success:    success,
		Error:      err.Error(),
	}
}
