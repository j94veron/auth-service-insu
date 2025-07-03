package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/j94veron/auth-service-insu/internal/models"
)

type TokenService struct {
	accessSecret  string
	refreshSecret string
}

type TokenClaims struct {
	jwt.StandardClaims
	UserID         uint   `json:"user_id"`
	Name           string `json:"name"`
	LastName       string `json:"last_name"`
	CommercialZone string `json:"commercial_zone"`
	Warehouse      string `json:"warehouse"`
	RoleID         uint   `json:"role_id"`
	OtherWarehouse string `json:"other_warehouse"`
	Province       string `json:"province"`
	Reports        string `json:"reports"`
	TokenUuid      string `json:"token_uuid"`
}

// NewTokenService creates a new instance of TokenService
func NewTokenService(accessSecret, refreshSecret string) *TokenService {
	return &TokenService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

// GenerateTokens creates new access and refresh tokens for a user
func (t *TokenService) GenerateTokens(user *models.User) (*models.TokenDetail, error) {
	return t.CreateTokens(user)
}

// CreateTokens creates the actual tokens with claims
func (t *TokenService) CreateTokens(user *models.User) (*models.TokenDetail, error) {
	td := &models.TokenDetail{}
	now := time.Now()

	// Configure expiration times
	td.AtExpires = now.Add(15 * time.Minute) // Access token: 15 minutes
	td.RtExpires = now.Add(2 * time.Hour)    // Refresh token: 2 hours

	// Generate UUIDs for tokens
	td.AccessUuid = uuid.New().String()
	td.RefreshUuid = uuid.New().String()

	// Create access token
	atClaims := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: td.AtExpires.Unix(),
			IssuedAt:  now.Unix(),
		},
		UserID:         user.ID,
		Name:           user.Name,
		LastName:       user.LastName,
		CommercialZone: user.CommercialZone,
		Warehouse:      user.Warehouse,
		OtherWarehouse: user.OtherWarehouse,
		Province:       user.Province,
		RoleID:         user.RoleID,
		Reports:        user.Reports,
		TokenUuid:      td.AccessUuid,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte(t.accessSecret))
	if err != nil {
		return nil, err
	}

	// Create refresh token
	rtClaims := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: td.RtExpires.Unix(),
			IssuedAt:  now.Unix(),
		},
		UserID:    user.ID,
		TokenUuid: td.RefreshUuid,
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(t.refreshSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

// VerifyToken checks if a token is valid
func (t *TokenService) VerifyToken(tokenString string, isRefresh bool) (*TokenClaims, error) {
	secret := t.accessSecret
	if isRefresh {
		secret = t.refreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken checks if the refresh token is valid
func (t *TokenService) ValidateRefreshToken(tokenString string) (claims *TokenClaims, err error) {
	claims = &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.refreshSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}
