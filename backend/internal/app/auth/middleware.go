package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"net/http"
	"strings"
)

var (
	ErrKIDNotFound         = fmt.Errorf("kid not found in token")
	ErrKIDInKeySetNotFound = fmt.Errorf("kid not found in key set")
	ErrInvalidToken        = fmt.Errorf("invalid token")
)

type CustomClaims struct {
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`

	jwt.RegisteredClaims
}

const (
	authorizationHeader = "Authorization"
	prefixTokenBearer   = "Bearer "
)

func AuthMiddleware(jwkSet jwk.Set) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "auth header required")
			return
		}

		parsedClaims, err := validateToken(jwkSet, strings.TrimPrefix(authHeader, prefixTokenBearer))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}

		c.Set("roles", parsedClaims.RealmAccess.Roles)
		c.Next()
	}
}
func validateToken(keySet jwk.Set, token string) (*CustomClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, ErrKIDNotFound
		}

		key, found := keySet.LookupKeyID(kid)
		if !found {
			return nil, ErrKIDInKeySetNotFound
		}

		var publicKey interface{}
		if err := key.Raw(&publicKey); err != nil {
			return nil, fmt.Errorf("key.Raw: %w", err)
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*CustomClaims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
