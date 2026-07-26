package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/shared"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		if authorizationHeader == "" {
			abortUnauthorized(c, "token não informado")
			return
		}

		tokenString, err := extractBearerToken(authorizationHeader)

		if err != nil {
			abortUnauthorized(c, err.Error())
			return
		}

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New(
						"método de assinatura inválido",
					)
				}

				return []byte(jwtSecret), nil
			},
			jwt.WithValidMethods([]string{
				jwt.SigningMethodHS256.Alg(),
			}),
		)

		if err != nil || !token.Valid {
			abortUnauthorized(c, "token inválido ou expirado")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			abortUnauthorized(c, "claims do token inválidas")
			return
		}

		userID, ok := claims["sub"].(string)

		if !ok || userID == "" {
			abortUnauthorized(
				c,
				"identificador do usuário inválido",
			)
			return
		}

		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		c.Set(shared.ContextUserIDKey, userID)
		c.Set(shared.ContextUserEmailKey, email)
		c.Set(shared.ContextUserRoleKey, role)

		c.Next()
	}
}

func extractBearerToken(header string) (string, error) {
	parts := strings.Fields(header)

	if len(parts) != 2 {
		return "", errors.New("formato do token inválido")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New(
			"tipo de autenticação inválido",
		)
	}

	if parts[1] == "" {
		return "", errors.New("token não informado")
	}

	return parts[1], nil
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		gin.H{
			"error": message,
		},
	)
}
