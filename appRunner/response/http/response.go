package zenHttpResponse

import (
	"errors"

	"github.com/labstack/echo/v4"
)

const ApiResponse = "api-response"
const ApiError = "api-error"

type Response interface {
	Val() interface{}
}

func NewJSONResponse(
	eCtx echo.Context,
	statusCode int,
	data interface{},
) error {
	var response Response

	if data == nil {
		return eCtx.JSON(statusCode, nil)
	}

	if statusCode >= 400 {
		err := data.(error)
		response = newErrorResponse(eCtx.Request().Context(), statusCode, err, nil)
		eCtx.Set(ApiError, err)
	} else {
		response = newSuccessResponse(data)
	}

	eCtx.Set(ApiResponse, response.Val())
	return eCtx.JSON(statusCode, response.Val())
}

func NewJSONCustomErrorResponse(
	eCtx echo.Context,
	statusCode int,
	err Error,
) error {
	eCtx.Set(ApiError, errors.New(err.Message))
	errorResponse := &errorResponse{Error: &err}
	eCtx.Set(ApiResponse, errorResponse.Val())
	return eCtx.JSON(statusCode, errorResponse.Val())
}

func NewJSONCustomErrorLocalizedResponse(
	eCtx echo.Context,
	statusCode int,
	data interface{},
	errLocalization map[string]map[string]string,
) error {
	var response Response

	if data == nil {
		return eCtx.JSON(statusCode, nil)
	}

	if statusCode >= 400 {
		err := data.(error)
		response = newErrorResponse(eCtx.Request().Context(), statusCode, err, errLocalization)
		eCtx.Set(ApiError, err)
	} else {
		response = newSuccessResponse(data)
	}

	eCtx.Set(ApiResponse, response.Val())
	return eCtx.JSON(statusCode, response.Val())
}
