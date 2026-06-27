package httpresponse

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func errorBody(message string, errDetail interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"message": message,
		"errors": errDetail,
	}
}

func Success(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func SuccessMessage(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]interface{}{
		"success": true,
		"message": message,
	})
}

func BadRequest(c echo.Context, message string, errDetail interface{}) error {
	return c.JSON(http.StatusBadRequest, errorBody(message, errDetail))
}

func ValidationFailed(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, errorBody("Validation failed", err.Error()))
}

func Unauthorized(c echo.Context, message string, errDetail interface{}) error {
	return c.JSON(http.StatusUnauthorized, errorBody(message, errDetail))
}

func Forbidden(c echo.Context, message string, errDetail interface{}) error {
	return c.JSON(http.StatusForbidden, errorBody(message, errDetail))
}

func NotFound(c echo.Context, message string, errDetail interface{}) error {
	return c.JSON(http.StatusNotFound, errorBody(message, errDetail))
}

func Conflict(c echo.Context, message string, errDetail interface{}) error {
	return c.JSON(http.StatusConflict, errorBody(message, errDetail))
}

func InternalError(c echo.Context, message string, errDetail interface{}) error {
	return c.JSON(http.StatusInternalServerError, errorBody(message, errDetail))
}

func HandleServiceError(c echo.Context, err error, notFoundMsg, defaultMsg string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(c, notFoundMsg, nil)
	}

	switch err.Error() {
	case "zone not found", "reservation not found":
		return NotFound(c, notFoundMsg, err.Error())
	case "forbidden: you can only cancel your own reservations":
		return Forbidden(c, "Forbidden", err.Error())
	case "email already registered":
		return BadRequest(c, err.Error(), err.Error())
	case "invalid email or password":
		return Unauthorized(c, err.Error(), err.Error())
	case "zone is at full capacity":
		return Conflict(c, "Zone is at full capacity", err.Error())
	}

	return InternalError(c, defaultMsg, err.Error())
}
