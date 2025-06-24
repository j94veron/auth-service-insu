package auth

import (
	"github.com/j94veron/auth-service-insu/logger"
	"net/http"
	"os"
	"time"

	"github.com/j94veron/auth-service-insu/internal/models"
	"github.com/j94veron/auth-service-insu/pkg/token"
)

// RefreshToken handles the request to refresh the token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Create a TokenService instance with the appropriate secrets
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	tokenService := token.NewTokenService(accessSecret, refreshSecret)

	refreshToken := r.Header.Get("Authorization")
	if refreshToken == "" {
		logger.Logger.Error("Refresh token required.")
		http.Error(w, "Missing refresh token", http.StatusBadRequest)
		return
	}

	claims, err := tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		logger.Logger.Error("Invalid Token " + err.Error())
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if time.Now().Unix() > claims.ExpiresAt {
		logger.Logger.Error("Refresh token expired")
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	user := models.User{
		ID:             claims.UserID,
		Name:           claims.Name,
		LastName:       claims.LastName,
		CommercialZone: claims.CommercialZone,
		Warehouse:      claims.Warehouse,
		OtherWarehouse: claims.OtherWarehouse,
		Province:       claims.Province,
	}

	tokens, err := tokenService.GenerateTokens(&user) // Generar nuevos tokens
	if err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"access_token":"` + tokens.AccessToken + `", "refresh_token":"` + tokens.RefreshToken + `"}`))
}
