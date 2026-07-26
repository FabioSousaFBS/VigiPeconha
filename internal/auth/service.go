package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, input RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, input LoginRequest) (*LoginResponse, error)
}

type AuthService struct {
	repository Repository
	jwtSecret  string
}

func NewService(repository Repository, jwtSecret string) Service {
	return &AuthService{
		repository: repository,
		jwtSecret:  jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterRequest) (*RegisterResponse, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if input.Name == "" {
		return nil, errors.New("nome é obrigatório")
	}

	if input.Email == "" {
		return nil, errors.New("email é obrigatório")
	}

	if len(input.Password) < 6 {
		return nil, errors.New("senha deve conter pelo menos 6 caracteres")
	}

	existingUser, err := s.repository.FindByEmail(ctx, input.Email)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email já cadastrado")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := users.User{
		ID:           uuid.NewString(),
		Name:         input.Name,
		Email:        input.Email,
		Phone:        "",
		PasswordHash: string(passwordHash),
		Role:         "citizen",
	}

	err = s.repository.CreateUser(ctx, user)

	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	input LoginRequest,
) (*LoginResponse, error) {

	user, err := s.repository.FindByEmail(
		ctx,
		strings.ToLower(strings.TrimSpace(input.Email)),
	)

	if err != nil {
		return nil, errors.New("usuário ou senha inválidos")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)

	if err != nil {
		return nil, errors.New("usuário ou senha inválidos")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(
		[]byte(s.jwtSecret),
	)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: tokenString,
		User: UserLoginPayload{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
