package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Data     any    `json:"data"`
	HttpCode int
}

func (resp *JsonResponse) ToJson(ctx *gin.Context) {
	code := 200
	if resp.HttpCode != 200 {
		code = resp.HttpCode
	}
	ctx.JSON(code, resp)
}

func FailResponse(code int, message string, data ...any) *JsonResponse {
	var r any
	if len(data) > 0 {
		r = data
	} else {
		r = struct{}{}
	}
	return &JsonResponse{
		Code:    code,
		Message: message,
		Data:    r,
	}
}

func SuccessResponse(data ...any) *JsonResponse {
	var r any
	if len(data) > 0 {
		r = data[0]
	} else {
		r = struct{}{}
	}
	return &JsonResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    r,
	}
}

func ErrorResponse(status int, message string, data ...any) *JsonResponse {
	var r any
	if len(data) > 0 {
		r = data
	} else {
		r = struct{}{}
	}

	return &JsonResponse{
		Code:    status,
		Message: message,
		Data:    r,
	}
}

func (resp *JsonResponse) WriteTo(ctx *gin.Context) {
	var code int
	if resp.HttpCode == 0 {
		code = http.StatusOK
	} else {
		code = resp.HttpCode
	}

	ctx.JSON(code, resp)
}
func (resp *JsonResponse) SetHttpCode(httpCode int) *JsonResponse {
	resp.HttpCode = httpCode
	return resp
}

func (that *JsonResponse) responseCode() int {
	if that.Code != http.StatusOK {
		return 200
	}
	return 200
}
