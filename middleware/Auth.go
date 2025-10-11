package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

//Remember to add Authorization header in your requests to
// protected routes, e.g. /api/user/me

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// 2. Must start with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// 3. Extract the token string (remove "Bearer ")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. Load secret key
		secret := []byte(os.Getenv("JWT_SECRET"))
		if len(secret) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not set"})
			c.Abort()
			return
		}

		// 5. Parse and verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})

		// 6. Handle invalid token cases
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or malformed token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Expired or invalid token"})
			c.Abort()
			return
		}

		// If you had multiple middleware layers, you might store user info in context
		// for downstream handlers to access.

		// 7. Extract claims (payload data)
		// claims, ok := token.Claims.(jwt.MapClaims)
		// if !ok {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		// 	c.Abort()
		// 	return
		// }

		// c.Set("user_id", claims["user_id"])
		// c.Set("email", claims["email"])

		c.Next()
	}
}
