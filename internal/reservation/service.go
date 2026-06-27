package reservation

import (
	"errors"

	"github.com/afifsylhet/spotsync-api/internal/reservation/dto"
	userdto "github.com/afifsylhet/spotsync-api/internal/user/dto"
	zonedto "github.com/afifsylhet/spotsync-api/internal/zone/dto"
	"github.com/afifsylhet/spotsync-api/internal/zone"
	"gorm.io/gorm"
)

type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	CancelReservation(reservationID, userID uint, userRole string) error
	GetAllReservations() ([]dto.AdminReservationResponse, error)
}

type reservationService struct {
	reservationRepo ReservationRepository
	zoneRepo        zone.ZoneRepository
}

func NewReservationService(
	reservationRepo ReservationRepository,
	zoneRepo zone.ZoneRepository,
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

	_, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("zone not found")
		}
		return nil, err
	}

	reservation, err := s.reservationRepo.CreateWithTransaction(
		userID,
		req.ZoneID,
		req.LicensePlate,
	)
	if err != nil {
		return nil, err
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
			Zone: zonedto.ZoneBasic{
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
			User: userdto.UserBasic{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			Zone: zonedto.ZoneBasic{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}

	return responses, nil
}
