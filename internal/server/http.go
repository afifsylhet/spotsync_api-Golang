package server

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/afifsylhet/spotsync-api/internal/auth"
	"github.com/afifsylhet/spotsync-api/internal/config"
	"github.com/afifsylhet/spotsync-api/internal/reservation"
	"github.com/afifsylhet/spotsync-api/internal/user"
	"github.com/afifsylhet/spotsync-api/internal/zone"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewServer() *echo.Echo {
	config.LoadEnv()

	db := config.ConnectDatabase()

	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo)
	zoneRepo := zone.NewZoneRepository(db)
	reservationRepo := reservation.NewReservationRepository(db)

	authSvc := auth.NewAuthService(userSvc)
	zoneSvc := zone.NewZoneService(zoneRepo)
	reservationSvc := reservation.NewReservationService(reservationRepo, zoneRepo)

	authHandler := auth.NewAuthHandler(authSvc)
	zoneHandler := zone.NewZoneHandler(zoneSvc)
	reservationHandler := reservation.NewReservationHandler(reservationSvc)

	e := echo.New()

	e.Validator = &CustomValidator{
		validator: validator.New(),
	}

	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
		},
	}))

	api := e.Group("/api/v1")

	auth.RegisterAuthRoutes(api, authHandler)
	zone.RegisterZoneRoutes(api, zoneHandler)
	reservation.RegisterReservationRoutes(api, reservationHandler)

	return e
}

func Start() {
	e := NewServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		if strings.Contains(err.Error(), "bind") || strings.Contains(err.Error(), "address already in use") {
			e.Logger.Fatal(fmt.Errorf(
				"port %s is already in use — stop the other server (Ctrl+C in its terminal) or set PORT to a different value",
				port,
			))
		}

		e.Logger.Fatal(err)
	}

	e.Listener = ln

	fmt.Printf("🚀 Server running on port %s\n", port)
	e.Logger.Fatal(e.Start(""))
}
