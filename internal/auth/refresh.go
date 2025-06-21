package auth

import (
	"net/http"
	"time"

	"github.com/j94veron/auth-service-insu/internal/models"
	"github.com/j94veron/auth-service-insu/pkg/token"
)

// RefreshToken handles the request to refresh the token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Crea una instancia de TokenService con los secretos apropiados
	tokenService := token.NewTokenService("your_access_secret", "your_refresh_secret") // Sustituye por tus secretos reales

	refreshToken := r.Header.Get("Authorization")
	if refreshToken == "" {
		http.Error(w, "Missing refresh token", http.StatusBadRequest)
		return
	}

	claims, err := tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if time.Now().Unix() > claims.ExpiresAt {
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
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"access_token":"` + tokens.AccessToken + `", "refresh_token":"` + tokens.RefreshToken + `"}`))
}
