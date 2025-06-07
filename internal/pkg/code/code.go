package code

import "github.com/cbhcbhcbh/Quantum/pkg/response"

var (
	OK = &response.JsonResponse{Code: 200, Message: "", Data: ""}

	InternalServerError = &response.JsonResponse{Code: 500, Message: "InternalError", Data: "Internal server error."}

	ErrPageNotFound = &response.JsonResponse{Code: 404, Message: "ResourceNotFound.PageNotFound", Data: "Page not found."}

	ErrBind = &response.JsonResponse{Code: 400, Message: "InvalidParameter.BindError", Data: "Error occurred while binding the request body to the struct."}

	ErrInvalidParameter = &response.JsonResponse{Code: 400, Message: "InvalidParameter", Data: "Parameter verification failed."}

	ErrSignToken = &response.JsonResponse{Code: 401, Message: "AuthFailure.SignTokenError", Data: "Error occurred while signing the JSON web token."}

	ErrTokenInvalid = &response.JsonResponse{Code: 401, Message: "AuthFailure.TokenInvalid", Data: "Token was invalid."}

	ErrUnauthorized = &response.JsonResponse{Code: 401, Message: "AuthFailure.Unauthorized", Data: "Unauthorized."}

	ErrUserAlreadyExist = &response.JsonResponse{Code: 400, Message: "FailedOperation.UserAlreadyExist", Data: "User already exist."}

	ErrEmailAlreadyExist = &response.JsonResponse{Code: 400, Message: "FailedOperation.EmailAlreadyExist", Data: "Email already exist."}

	ErrUserNotFound = &response.JsonResponse{Code: 404, Message: "ResourceNotFound.UserNotFound", Data: "User was not found."}

	ErrPasswordIncorrect = &response.JsonResponse{Code: 401, Message: "InvalidParameter.PasswordIncorrect", Data: "Password was incorrect."}
)
