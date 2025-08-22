package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saurabhraut1212/notes_sharing_api/internal/models"
	"github.com/saurabhraut1212/notes_sharing_api/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo  *repo.UserRepo
	JWTSecret string
}

func NewAuthHandler(userRepo *repo.UserRepo, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "email and password required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existing, _ := h.UserRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return c.Status(400).JSON(fiber.Map{"error": "email already registered"})
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	}

	if err := h.UserRepo.Create(ctx, u); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create user"})
	}
	// return basic user info
	return c.Status(201).JSON(fiber.Map{"message": "user created"})

}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, _ := token.SignedString([]byte(h.JWTSecret))
	return c.JSON(fiber.Map{"token": tokenStr})
}
