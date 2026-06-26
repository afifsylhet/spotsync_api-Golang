package service

import (
	"github.com/afifsylhet/spotsync-api/dto"
	"github.com/afifsylhet/spotsync-api/models"
	"github.com/afifsylhet/spotsync-api/repository"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
	UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	DeleteZone(id uint) error
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{
		zoneRepo: zoneRepo,
	}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}

	return toZoneResponse(zone, 0), nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ZoneResponse

	for _, zone := range zones {
		// Calculate available spots for each zone
		activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
		if err != nil {
			return nil, err
		}

		available := zone.TotalCapacity - int(activeCount)

		responses = append(responses, *toZoneResponse(&zone, available))
	}

	return responses, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, err
	}

	available := zone.TotalCapacity - int(activeCount)

	return toZoneResponse(zone, available), nil
}

func (s *zoneService) UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Only update fields that were provided
	if req.Name != "" {
		zone.Name = req.Name
	}

	if req.Type != "" {
		zone.Type = req.Type
	}

	if req.TotalCapacity > 0 {
		zone.TotalCapacity = req.TotalCapacity
	}

	if req.PricePerHour > 0 {
		zone.PricePerHour = req.PricePerHour
	}

	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, err
	}

	return toZoneResponse(zone, 0), nil
}

func (s *zoneService) DeleteZone(id uint) error {
	_, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.zoneRepo.Delete(id)
}

// Helper to convert model -> DTO
func toZoneResponse(zone *models.ParkingZone, available int) *dto.ZoneResponse {
	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}
}
