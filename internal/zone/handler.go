package zone

import (
	"net/http"
	"strconv"

	"github.com/afifsylhet/spotsync-api/internal/httpresponse"
	"github.com/afifsylhet/spotsync-api/internal/zone/dto"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ZoneHandler struct {
	zoneService ZoneService
}

func NewZoneHandler(zoneService ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zoneService}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest

	if err := c.Bind(&req); err != nil {
		return httpresponse.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := c.Validate(req); err != nil {
		return httpresponse.ValidationFailed(c, err)
	}

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		return httpresponse.InternalError(c, "Failed to create parking zone", err.Error())
	}

	return httpresponse.Success(c, http.StatusCreated, "Parking zone created successfully", zone)
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return httpresponse.InternalError(c, "Failed to retrieve parking zones", err.Error())
	}

	return httpresponse.Success(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return httpresponse.BadRequest(c, "Invalid zone ID", err.Error())
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return httpresponse.NotFound(c, "Zone not found", nil)
		}

		return httpresponse.InternalError(c, "Failed to retrieve parking zone", err.Error())
	}

	return httpresponse.Success(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}

func (h *ZoneHandler) UpdateZone(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return httpresponse.BadRequest(c, "Invalid zone ID", err.Error())
	}

	var req dto.UpdateZoneRequest

	if err := c.Bind(&req); err != nil {
		return httpresponse.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := c.Validate(req); err != nil {
		return httpresponse.ValidationFailed(c, err)
	}

	zone, err := h.zoneService.UpdateZone(uint(id), req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return httpresponse.NotFound(c, "Zone not found", nil)
		}

		return httpresponse.InternalError(c, "Failed to update parking zone", err.Error())
	}

	return httpresponse.Success(c, http.StatusOK, "Parking zone updated successfully", zone)
}

func (h *ZoneHandler) DeleteZone(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return httpresponse.BadRequest(c, "Invalid zone ID", err.Error())
	}

	if err := h.zoneService.DeleteZone(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return httpresponse.NotFound(c, "Zone not found", nil)
		}

		return httpresponse.InternalError(c, "Failed to delete parking zone", err.Error())
	}

	return httpresponse.SuccessMessage(c, http.StatusOK, "Parking zone deleted successfully")
}
