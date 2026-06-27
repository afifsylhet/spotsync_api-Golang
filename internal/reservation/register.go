package reservation

import (
	"github.com/afifsylhet/spotsync-api/internal/auth"
	"github.com/labstack/echo/v4"
)

func RegisterReservationRoutes(g *echo.Group, handler *ReservationHandler) {
	reservations := g.Group("/reservations")
	reservations.Use(auth.JWTMiddleware())

	reservations.POST("", handler.CreateReservation)
	reservations.GET("/my-reservations", handler.GetMyReservations)
	reservations.DELETE("/:id", handler.CancelReservation)

	reservations.GET("",
		handler.GetAllReservations,
		auth.AdminOnly,
	)
}
