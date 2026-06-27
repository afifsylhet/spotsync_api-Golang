package auth

import (
	"net/http"

	"github.com/afifsylhet/spotsync-api/internal/auth/dto"
	"github.com/afifsylhet/spotsync-api/internal/httpresponse"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return httpresponse.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := c.Validate(req); err != nil {
		return httpresponse.ValidationFailed(c, err)
	}

	user, err := h.authService.Register(req)
	if err != nil {
		return httpresponse.BadRequest(c, err.Error(), err.Error())
	}

	return httpresponse.Success(c, http.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return httpresponse.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := c.Validate(req); err != nil {
		return httpresponse.ValidationFailed(c, err)
	}

	result, err := h.authService.Login(req)
	if err != nil {
		return httpresponse.Unauthorized(c, err.Error(), err.Error())
	}

	return httpresponse.Success(c, http.StatusOK, "Login successful", result)
}
