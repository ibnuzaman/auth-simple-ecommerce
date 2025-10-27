package helpers

import "github.com/labstack/echo/v4"

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseHttp(e echo.Context, code int, message string, data interface{}) error {
	resp := BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}

	return e.JSON(code, resp)
}
