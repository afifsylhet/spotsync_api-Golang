package reservation

import (
	"net/http"
	"strconv"

	"github.com/afifsylhet/spotsync-api/internal/auth"
	"github.com/afifsylhet/spotsync-api/internal/httpresponse"
	"github.com/afifsylhet/spotsync-api/internal/reservation/dto"
	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	reservationService ReservationService
}

func NewReservationHandler(reservationService ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
	}
}

func (h *ReservationHandler) CreateReservation(c echo.Context) error {
	userID, _ := auth.GetUserClaims(c)

	var req dto.CreateReservationRequest

	if err := c.Bind(&req); err != nil {
		return httpresponse.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := c.Validate(req); err != nil {
		return httpresponse.ValidationFailed(c, err)
	}

	reservation, err := h.reservationService.CreateReservation(userID, req)
	if err != nil {
		if err == ErrZoneFull {
			return httpresponse.Conflict(c, "Zone is at full capacity", err.Error())
		}

		if err.Error() == "zone not found" {
			return httpresponse.NotFound(c, "Zone not found", err.Error())
		}

		return httpresponse.InternalError(c, "Failed to create reservation", err.Error())
	}

	return httpresponse.Success(c, http.StatusCreated, "Reservation confirmed successfully", reservation)
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, _ := auth.GetUserClaims(c)

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return httpresponse.InternalError(c, "Failed to retrieve reservations", err.Error())
	}

	return httpresponse.Success(c, http.StatusOK, "My reservations retrieved successfully", reservations)
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	userID, userRole := auth.GetUserClaims(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return httpresponse.BadRequest(c, "Invalid reservation ID", err.Error())
	}

	if err := h.reservationService.CancelReservation(uint(id), userID, userRole); err != nil {
		if err.Error() == "reservation not found" {
			return httpresponse.NotFound(c, "Reservation not found", err.Error())
		}

		if err.Error() == "forbidden: you can only cancel your own reservations" {
			return httpresponse.Forbidden(c, "Forbidden", err.Error())
		}

		return httpresponse.BadRequest(c, "Failed to cancel reservation", err.Error())
	}

	return httpresponse.SuccessMessage(c, http.StatusOK, "Reservation cancelled successfully")
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return httpresponse.InternalError(c, "Failed to retrieve reservations", err.Error())
	}

	return httpresponse.Success(c, http.StatusOK, "All reservations retrieved successfully", reservations)
}
