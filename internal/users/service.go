package users

import (
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetLicenseUsage(
		ctx context.Context,
		organizationID string,
	) (*LicenseUsage, error)

	Create(
		ctx context.Context,
		organizationID string,
		request CreateUserRequest,
	) (*User, error)

	UpdateStatus(
		ctx context.Context,
		userID string,
		organizationID string,
		status string,
	) (*User, error)
}

type UserService struct {
	repository Repository
}

func NewService(
	repository Repository,
) Service {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) GetLicenseUsage(
	ctx context.Context,
	organizationID string,
) (*LicenseUsage, error) {
	return s.repository.GetLicenseUsage(
		ctx,
		organizationID,
	)
}

func (s *UserService) Create(
	ctx context.Context,
	organizationID string,
	request CreateUserRequest,
) (*User, error) {
	name := strings.TrimSpace(request.Name)
	email := strings.ToLower(
		strings.TrimSpace(request.Email),
	)

	if name == "" {
		return nil, ErrInvalidName
	}

	if email == "" {
		return nil, ErrInvalidEmail
	}

	if len(request.Password) < 8 {
		return nil, ErrInvalidPassword
	}

	if !isAllowedOrganizationRole(request.Role) {
		return nil, ErrInvalidRole
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	var phone *string

	if request.Phone != nil {
		value := strings.TrimSpace(*request.Phone)

		if value != "" {
			phone = &value
		}
	}

	return s.repository.CreateWithLicenseValidation(
		ctx,
		CreateUserParams{
			Name:           name,
			Email:          email,
			Phone:          phone,
			PasswordHash:   string(passwordHash),
			Role:           request.Role,
			OrganizationID: organizationID,
		},
	)
}

func (s *UserService) UpdateStatus(
	ctx context.Context,
	userID string,
	organizationID string,
	status string,
) (*User, error) {
	if status != "active" &&
		status != "inactive" {
		return nil, ErrInvalidStatus
	}

	return s.repository.
		UpdateStatusWithLicenseValidation(
			ctx,
			userID,
			organizationID,
			status,
		)
}

func isAllowedOrganizationRole(
	role string,
) bool {
	switch role {
	case "organization_admin",
		"agent",
		"analyst",
		"viewer":
		return true

	default:
		return false
	}
}
