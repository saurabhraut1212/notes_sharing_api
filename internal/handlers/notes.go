package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/saurabhraut1212/notes_sharing_api/internal/config"
	"github.com/saurabhraut1212/notes_sharing_api/internal/models"
	"github.com/saurabhraut1212/notes_sharing_api/internal/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteHandler struct {
	NoteRepo *repo.NoteRepo
	Config   *config.Config
}

func NewNoteHandler(noteRepo *repo.NoteRepo, cfg *config.Config) *NoteHandler {
	return &NoteHandler{
		NoteRepo: noteRepo,
		Config:   cfg,
	}
}

func (h *NoteHandler) CreateNote(c *fiber.Ctx) error {
	var req struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		IsPublic bool     `json:"is_public"`
		Tags     []string `json:"tags"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	// get user id from locals
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userId := uid.(primitive.ObjectID)

	n := &models.Note{
		UserID:   userId,
		Title:    req.Title,
		Content:  req.Content,
		IsPublic: req.IsPublic,
		Tags:     req.Tags,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.NoteRepo.Create(ctx, n); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create note"})
	}
	return c.Status(201).JSON(n)

}

func (h *NoteHandler) GetMyNotes(c *fiber.Ctx) error {
	userIDIface := c.Locals("user_id")
	if userIDIface == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := userIDIface.(primitive.ObjectID)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := h.NoteRepo.ListByUser(ctx, userID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch notes"})
	}
	return c.JSON(fiber.Map{"notes": items})
}

func (h *NoteHandler) GetPublicNotes(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	items, err := h.NoteRepo.ListPublic(ctx, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func (h *NoteHandler) GetNoteByID(c *fiber.Ctx) error {
	idHex := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	n, err := h.NoteRepo.GetById(ctx, oid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if n == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	// if note is private, ensure owner or authenticated user
	if !n.IsPublic {
		userIDIface := c.Locals("user_id")
		if userIDIface == nil || userIDIface.(primitive.ObjectID) != n.UserID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}
	return c.JSON(n)
}

func (h *NoteHandler) UpdateNote(c *fiber.Ctx) error {
	idHex := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	// ensure owner
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := h.NoteRepo.GetById(ctx, oid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if n == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	userIDIface := c.Locals("user_id")
	if userIDIface == nil || userIDIface.(primitive.ObjectID) != n.UserID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	update := bson.M{}
	if v, ok := req["title"].(string); ok {
		update["title"] = v
	}
	if v, ok := req["content"].(string); ok {
		update["content"] = v
	}
	if v, ok := req["is_public"].(bool); ok {
		update["is_public"] = v
	}
	if v, ok := req["tags"].([]interface{}); ok {
		var tags []string
		for _, iv := range v {
			if s, ok := iv.(string); ok {
				tags = append(tags, s)
			}
		}
		update["tags"] = tags
	}

	updated, err := h.NoteRepo.Update(ctx, oid, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if updated == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(updated)
}

func (h *NoteHandler) DeleteNote(c *fiber.Ctx) error {
	idHex := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	// ensure owner
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := h.NoteRepo.GetById(ctx, oid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if n == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	userIDIface := c.Locals("user_id")
	if userIDIface == nil || userIDIface.(primitive.ObjectID) != n.UserID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if err := h.NoteRepo.Delete(ctx, oid); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "note deleted"})
}
