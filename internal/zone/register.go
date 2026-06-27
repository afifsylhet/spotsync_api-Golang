package zone

import (
	"github.com/afifsylhet/spotsync-api/internal/auth"
	"github.com/labstack/echo/v4"
)

func RegisterZoneRoutes(g *echo.Group, handler *ZoneHandler) {
	zones := g.Group("/zones")
	zones.GET("", handler.GetAllZones)
	zones.GET("/:id", handler.GetZoneByID)

	zones.POST("",
		handler.CreateZone,
		auth.JWTMiddleware(),
		auth.AdminOnly,
	)

	zones.PUT("/:id",
		handler.UpdateZone,
		auth.JWTMiddleware(),
		auth.AdminOnly,
	)

	zones.DELETE("/:id",
		handler.DeleteZone,
		auth.JWTMiddleware(),
		auth.AdminOnly,
	)
}
