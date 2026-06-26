package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/afifsylhet/spotsync-api/config"
	"github.com/afifsylhet/spotsync-api/handler"
	custommiddleware "github.com/afifsylhet/spotsync-api/middleware"
	"github.com/afifsylhet/spotsync-api/repository"
	"github.com/afifsylhet/spotsync-api/service"
)

// CustomValidator wraps go-playground/validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	// Connect to database
	db := config.ConnectDatabase()

	// ── DEPENDENCY INJECTION ──────────────────────────────────────

	// 1. Repositories (talk to DB)
	userRepo := repository.NewUserRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	// 2. Services (talk to Repositories)
	authSvc := service.NewAuthService(userRepo)
	zoneSvc := service.NewZoneService(zoneRepo)
	reservationSvc := service.NewReservationService(reservationRepo, zoneRepo)

	// 3. Handlers (talk to Services)
	authHandler := handler.NewAuthHandler(authSvc)
	zoneHandler := handler.NewZoneHandler(zoneSvc)
	reservationHandler := handler.NewReservationHandler(reservationSvc)

	// ─────────────────────────────────────────────────────────────

	// Setup Echo
	e := echo.New()

	// Attach validator
	e.Validator = &CustomValidator{
		validator: validator.New(),
	}

	// Global Middleware
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

	// ── ROUTES ───────────────────────────────────────────────────

	api := e.Group("/api/v1")

	// Auth (Public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Zones (Mixed access)
	zones := api.Group("/zones")
	zones.GET("", zoneHandler.GetAllZones)     // Public
	zones.GET("/:id", zoneHandler.GetZoneByID) // Public

	zones.POST("",
		zoneHandler.CreateZone,
		custommiddleware.JWTMiddleware(),
		custommiddleware.AdminOnly,
	)

	zones.PUT("/:id",
		zoneHandler.UpdateZone,
		custommiddleware.JWTMiddleware(),
		custommiddleware.AdminOnly,
	)

	zones.DELETE("/:id",
		zoneHandler.DeleteZone,
		custommiddleware.JWTMiddleware(),
		custommiddleware.AdminOnly,
	)

	// Reservations
	reservations := api.Group("/reservations")
	reservations.Use(custommiddleware.JWTMiddleware()) // ALL reservation routes

	reservations.POST("", reservationHandler.CreateReservation)
	reservations.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservations.DELETE("/:id", reservationHandler.CancelReservation)

	reservations.GET("",
		reservationHandler.GetAllReservations,
		custommiddleware.AdminOnly,
	) // Admin only

	// ─────────────────────────────────────────────────────────────

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server running on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
