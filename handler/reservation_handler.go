package handler

import (
	"net/http"
	"strconv"

	"github.com/afifsylhet/spotsync-api/dto"
	"github.com/afifsylhet/spotsync-api/repository"
	"github.com/afifsylhet/spotsync-api/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	reservationService service.ReservationService
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
	}
}

// Helper to extract user claims from JWT in Echo context
func getUserClaims(c echo.Context) (uint, string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	return userID, role
}

func (h *ReservationHandler) CreateReservation(c echo.Context) error {
	userID, _ := getUserClaims(c)

	var req dto.CreateReservationRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body", err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Validation failed", err.Error()))
	}

	reservation, err := h.reservationService.CreateReservation(userID, req)
	if err != nil {
		// Zone full = 409 Conflict
		if err == repository.ErrZoneFull {
			return c.JSON(http.StatusConflict, errorResponse("Zone is at full capacity", err.Error()))
		}

		if err.Error() == "zone not found" {
			return c.JSON(http.StatusNotFound, errorResponse("Zone not found", err.Error()))
		}

		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to create reservation", err.Error()))
	}

	return c.JSON(http.StatusCreated, successResponse("Reservation confirmed successfully", reservation))
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, _ := getUserClaims(c)

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to retrieve reservations", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("My reservations retrieved successfully", reservations))
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	userID, userRole := getUserClaims(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid reservation ID", err.Error()))
	}

	if err := h.reservationService.CancelReservation(uint(id), userID, userRole); err != nil {
		if err.Error() == "reservation not found" {
			return c.JSON(http.StatusNotFound, errorResponse("Reservation not found", err.Error()))
		}

		if err.Error() == "forbidden: you can only cancel your own reservations" {
			return c.JSON(http.StatusForbidden, errorResponse("Forbidden", err.Error()))
		}

		return c.JSON(http.StatusBadRequest, errorResponse("Failed to cancel reservation", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Reservation cancelled successfully", nil))
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to retrieve reservations", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("All reservations retrieved successfully", reservations))
}
