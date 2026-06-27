package auth

import "github.com/labstack/echo/v4"

func RegisterAuthRoutes(g *echo.Group, handler *AuthHandler) {
	auth := g.Group("/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
}
