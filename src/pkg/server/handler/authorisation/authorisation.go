package authorisation

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	domain "templify/pkg/domain/model"
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
		ah.HandleErrors(w, r, domain.ErrAuthorizationHeaderMissing)
		return false //unauthorized
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer")
	if tokenString == authHeader {
		ah.HandleErrors(w, r, domain.ErrInvalidTokenFormat)
		return false //unauthorized
	}
	token, err := a.VerifyToken(tokenString)
	if err != nil || !token.Valid {
		ah.HandleErrors(w, r, domain.ErrInvalidToken)
		return false //unauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ah.HandleErrors(w, r, domain.ErrAccessDenied)
		return false //unauthorized
	}
	for key, value := range requiredClaims {
		if claims[key] != value {
			ah.HandleErrors(w, r, domain.ErrAccessDenied)
			return false //unauthorized
		}
	}
	return true //authorized
}

func (a *Authorizer) VerifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return token, nil
	})
}
