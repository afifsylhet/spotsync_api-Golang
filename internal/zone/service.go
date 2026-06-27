package zone

import (
	"github.com/afifsylhet/spotsync-api/internal/zone/dto"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneCreateResponse, error)
	GetAllZones() ([]dto.ZoneListResponse, error)
	GetZoneByID(id uint) (*dto.ZoneListResponse, error)
	UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.ZoneCreateResponse, error)
	DeleteZone(id uint) error
}

type zoneService struct {
	zoneRepo ZoneRepository
}

func NewZoneService(zoneRepo ZoneRepository) ZoneService {
	return &zoneService{
		zoneRepo: zoneRepo,
	}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneCreateResponse, error) {
	zone := &ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}

	return toZoneCreateResponse(zone), nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneListResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ZoneListResponse

	for _, zone := range zones {
		activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
		if err != nil {
			return nil, err
		}

		available := zone.TotalCapacity - int(activeCount)

		responses = append(responses, *toZoneListResponse(&zone, available))
	}

	return responses, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneListResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, err
	}

	available := zone.TotalCapacity - int(activeCount)

	return toZoneListResponse(zone, available), nil
}

func (s *zoneService) UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.ZoneCreateResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

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

	return toZoneCreateResponse(zone), nil
}

func (s *zoneService) DeleteZone(id uint) error {
	_, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.zoneRepo.Delete(id)
}

func toZoneCreateResponse(zone *ParkingZone) *dto.ZoneCreateResponse {
	return &dto.ZoneCreateResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}
}

func toZoneListResponse(zone *ParkingZone, available int) *dto.ZoneListResponse {
	return &dto.ZoneListResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}
}
