package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/saurabhraut1212/notes_sharing_api/internal/repo"
)

type TagHandler struct {
	TagRepo *repo.TagRepo
}

func NewTagHandler(tagRepo *repo.TagRepo) *TagHandler {
	return &TagHandler{
		TagRepo: tagRepo,
	}
}

func (h *TagHandler) TopTags(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	arr, err := h.TagRepo.TopTags(ctx, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// convert to friendly JSON
	out := make([]fiber.Map, 0, len(arr))
	for _, a := range arr {
		out = append(out, fiber.Map{"tag": a["_id"], "count": a["count"]})
	}
	return c.JSON(out)
}
