package auth

import (
	"context"
	"errors"
	"time"

	"github.com/j94veron/auth-service-insu/internal/models"
	"github.com/j94veron/auth-service-insu/internal/user"
	"github.com/j94veron/auth-service-insu/pkg/redis"
	"github.com/j94veron/auth-service-insu/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo     user.Repository
	tokenService *token.TokenService
	redisClient  *redis.Client
}

func NewService(userRepo user.Repository, tokenService *token.TokenService, redisClient *redis.Client) *Service {
	return &Service{
		userRepo:     userRepo,
		tokenService: tokenService,
		redisClient:  redisClient,
	}
}

func (s *Service) Login(email, password, endpoint string) (*models.TokenDetail, *models.User, error) {
	// Check if the email exists
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, nil,
			errors.New("usuario no encontrado")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, errors.New("contraseña incorrecta")
	}

	// Generate token
	td, err := s.tokenService.CreateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	// Save token in Redis
	ctx := context.Background()
	if err := s.redisClient.SaveToken(ctx, td.AccessUuid, user.ID, time.Until(td.AtExpires)); err != nil {
		return nil, nil, err
	}

	if err := s.redisClient.SaveToken(ctx, td.RefreshUuid, user.ID, time.Until(td.RtExpires)); err != nil {
		return nil, nil, err
	}

	return td, user, nil
}

func (s *Service) Refresh(refreshToken string) (*models.TokenDetail, error) {

	// Check refresh token
	claims, err := s.tokenService.VerifyToken(refreshToken, true)
	if err != nil {
		return nil, errors.New("refresh token inválido")
	}

	// Check if the token is in Redis
	ctx := context.Background()
	_, err = s.redisClient.GetUserID(ctx, claims.TokenUuid)
	if err != nil {
		return nil, errors.New("refresh token revocado o expirado")
	}

	// Search user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Generate new tokens
	td, err := s.tokenService.CreateTokens(user)
	if err != nil {
		return nil, err
	}

	// Delete the old refresh token and save the new ones
	if err := s.redisClient.DeleteToken(ctx, claims.TokenUuid); err != nil {
		return nil, err
	}

	if err := s.redisClient.SaveToken(ctx, td.AccessUuid, user.ID, time.Until(td.AtExpires)); err != nil {
		return nil, err
	}

	if err := s.redisClient.SaveToken(ctx, td.RefreshUuid, user.ID, time.Until(td.RtExpires)); err != nil {
		return nil, err
	}

	return td, nil
}

func (s *Service) hasPermissionForEndpoint(user *models.User, endpoint string) bool {
	//Check if the user role has permission for the endpoint
	for _, perm := range user.Role.Permissions {
		if perm.Endpoint == endpoint {
			return true
		}
	}
	return false
}

func (s *Service) Logout(userID uint, accessUuid string) error {
	ctx := context.Background()
	return s.redisClient.DeleteToken(ctx, accessUuid)
}
