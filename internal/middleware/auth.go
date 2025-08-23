package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saurabhraut1212/notes_sharing_api/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RequireAuth verifies Bearer token and sets "user_id" local (primitive.ObjectID)
func RequireAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing authorization header"})
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "invalid authorization header"})
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token claims"})
		}
		uidStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid user id in token"})
		}
		oid, err := primitive.ObjectIDFromHex(uidStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid user id"})
		}
		// optional: check exp
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return c.Status(401).JSON(fiber.Map{"error": "token expired"})
			}
		}
		// set user id to locals for handlers
		c.Locals("user_id", oid)
		return c.Next()
	}
}
