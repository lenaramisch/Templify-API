package authorisation

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	ah "templify/pkg/server/handler"

	"github.com/golang-jwt/jwt"
)

type Authorizer struct {
	FeatureFlag bool
	log         *slog.Logger
}

func NewAuthorizer(featureFlag bool, logger *slog.Logger) *Authorizer {
	return &Authorizer{
		FeatureFlag: featureFlag,
		log:         logger,
	}
}

func (a *Authorizer) CheckIfAuthorised(w http.ResponseWriter, r *http.Request, requiredClaims map[string]any) bool {
	if !a.FeatureFlag {
		return true //authorized
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		ah.HandleError(w, r, http.StatusUnauthorized, "Authorization header missing")
		return false //unauthorized
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer")
	if tokenString == authHeader {
		ah.HandleError(w, r, http.StatusUnauthorized, "Invalid token format, expected Bearer")
		return false //unauthorized
	}
	token, err := a.VerifyToken(tokenString)
	if err != nil || !token.Valid {
		ah.HandleError(w, r, http.StatusUnauthorized, "Invalid token")
		return false //unauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ah.HandleError(w, r, http.StatusForbidden, "Access denied")
		return false //unauthorized
	}
	for key, value := range requiredClaims {
		if claims[key] != value {
			ah.HandleError(w, r, http.StatusForbidden, "Access denied")
			return false //unauthorized
		}
	}
	return true //authorized
}

func (a *Authorizer) VerifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return token, nil
	})
}
