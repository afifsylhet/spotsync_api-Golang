package handler

import (
	"net/http"
	"strconv"

	"github.com/afifsylhet/spotsync-api/dto"
	"github.com/afifsylhet/spotsync-api/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ZoneHandler struct {
	zoneService service.ZoneService
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zoneService}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body", err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Validation failed", err.Error()))
	}

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to create parking zone", err.Error()))
	}

	return c.JSON(http.StatusCreated, successResponse("Parking zone created successfully", zone))
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to retrieve parking zones", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Parking zones retrieved successfully", zones))
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid zone ID", err.Error()))
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, errorResponse("Zone not found", nil))
		}

		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to retrieve parking zone", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Parking zone retrieved successfully", zone))
}

func (h *ZoneHandler) UpdateZone(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid zone ID", err.Error()))
	}

	var req dto.UpdateZoneRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid request body", err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Validation failed", err.Error()))
	}

	zone, err := h.zoneService.UpdateZone(uint(id), req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, errorResponse("Zone not found", nil))
		}

		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to update parking zone", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Parking zone updated successfully", zone))
}

func (h *ZoneHandler) DeleteZone(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Invalid zone ID", err.Error()))
	}

	if err := h.zoneService.DeleteZone(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, errorResponse("Zone not found", nil))
		}

		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to delete parking zone", err.Error()))
	}

	return c.JSON(http.StatusOK, successResponse("Parking zone deleted successfully", nil))
}

// Helper functions to keep handlers clean

func successResponse(message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}
}

func errorResponse(message string, errors interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"message": message,
		"errors":  errors,
	}
}
