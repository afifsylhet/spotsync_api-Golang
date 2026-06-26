package service

import (
	"errors"

	"github.com/afifsylhet/spotsync-api/dto"
	"github.com/afifsylhet/spotsync-api/repository"
	"gorm.io/gorm"
)

type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	CancelReservation(reservationID, userID uint, userRole string) error
	GetAllReservations() ([]dto.AdminReservationResponse, error)
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	zoneRepo        repository.ZoneRepository
}

func NewReservationService(
	reservationRepo repository.ReservationRepository,
	zoneRepo repository.ZoneRepository,
) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

func (s *reservationService) CreateReservation(
	userID uint,
	req dto.CreateReservationRequest,
) (*dto.ReservationResponse, error) {

	// Verify zone exists first
	_, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("zone not found")
		}
		return nil, err
	}

	// Create reservation with transaction + row lock (concurrency safe)
	reservation, err := s.reservationRepo.CreateWithTransaction(
		userID,
		req.ZoneID,
		req.LicensePlate,
	)
	if err != nil {
		return nil, err // ErrZoneFull or db error
	}

	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}, nil
}

func (s *reservationService) GetMyReservations(
	userID uint,
) ([]dto.MyReservationResponse, error) {

	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.MyReservationResponse

	for _, r := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneBasic{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}

	return responses, nil
}

func (s *reservationService) CancelReservation(
	reservationID,
	userID uint,
	userRole string,
) error {

	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reservation not found")
		}
		return err
	}

	// Drivers can only cancel their OWN reservations
	if userRole != "admin" && reservation.UserID != userID {
		return errors.New("forbidden: you can only cancel your own reservations")
	}

	if reservation.Status == "cancelled" {
		return errors.New("reservation is already cancelled")
	}

	return s.reservationRepo.Cancel(reservationID)
}

func (s *reservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {

	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.AdminReservationResponse

	for _, r := range reservations {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User: dto.UserBasic{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			Zone: dto.ZoneBasic{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}

	return responses, nil
}
