package healthfacilities

import (
	"context"
	"math"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	Create(context.Context, SaveRequest) (*Facility, error)
	Update(context.Context, string, SaveRequest) (*Facility, error)
	GetByID(context.Context, string) (*Facility, error)
	List(context.Context, ListFilter) ([]Facility, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func validateRequest(r *SaveRequest) error {
	r.Name = strings.TrimSpace(r.Name)
	r.UnitType = strings.ToLower(strings.TrimSpace(r.UnitType))
	r.City = strings.TrimSpace(r.City)
	r.State = strings.ToUpper(strings.TrimSpace(r.State))
	if r.Name == "" || r.City == "" || len(r.State) != 2 || r.Latitude == nil || r.Longitude == nil {
		return ErrInvalidInput
	}
	if math.IsNaN(*r.Latitude) || math.IsNaN(*r.Longitude) || math.IsInf(*r.Latitude, 0) || math.IsInf(*r.Longitude, 0) || *r.Latitude < -90 || *r.Latitude > 90 || *r.Longitude < -180 || *r.Longitude > 180 {
		return ErrInvalidInput
	}
	switch r.UnitType {
	case "hospital", "upa", "emergency", "reference_center", "other":
	default:
		return ErrInvalidInput
	}
	if r.CNES != nil {
		v := strings.TrimSpace(*r.CNES)
		if v == "" {
			r.CNES = nil
		} else {
			r.CNES = &v
		}
	}
	if r.Phone != nil {
		v := strings.TrimSpace(*r.Phone)
		if v == "" {
			r.Phone = nil
		} else {
			r.Phone = &v
		}
	}
	if r.Address != nil {
		v := strings.TrimSpace(*r.Address)
		if v == "" {
			r.Address = nil
		} else {
			r.Address = &v
		}
	}
	return nil
}
func (s *service) Create(ctx context.Context, r SaveRequest) (*Facility, error) {
	if err := validateRequest(&r); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, r)
}
func (s *service) Update(ctx context.Context, id string, r SaveRequest) (*Facility, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	if err := validateRequest(&r); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, r)
}
func (s *service) GetByID(ctx context.Context, id string) (*Facility, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}
	return s.repo.GetByID(ctx, id)
}
func (s *service) List(ctx context.Context, f ListFilter) ([]Facility, error) {
	f.State = strings.ToUpper(strings.TrimSpace(f.State))
	f.City = strings.TrimSpace(f.City)
	if f.State != "" && len(f.State) != 2 {
		return nil, ErrInvalidInput
	}
	if (f.Latitude == nil) != (f.Longitude == nil) {
		return nil, ErrInvalidInput
	}
	if f.Latitude != nil {
		if math.IsNaN(*f.Latitude) || math.IsNaN(*f.Longitude) || math.IsInf(*f.Latitude, 0) || math.IsInf(*f.Longitude, 0) || math.IsNaN(f.RadiusMeters) || math.IsInf(f.RadiusMeters, 0) || *f.Latitude < -90 || *f.Latitude > 90 || *f.Longitude < -180 || *f.Longitude > 180 || f.RadiusMeters <= 0 || f.RadiusMeters > 200000 {
			return nil, ErrInvalidInput
		}
	}
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	return s.repo.List(ctx, f)
}
