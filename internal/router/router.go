package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/saurabhraut1212/notes_sharing_api/internal/config"
	"github.com/saurabhraut1212/notes_sharing_api/internal/handlers"
	"github.com/saurabhraut1212/notes_sharing_api/internal/middleware"
	"github.com/saurabhraut1212/notes_sharing_api/internal/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(client *mongo.Client, cfg *config.Config) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	//repos
	userRepo := repo.NewUserRepo(client.Database(cfg.DBName))
	noteRepo := repo.NewNoteRepo(client.Database(cfg.DBName))
	tagRepo := repo.NewTagRepo(client.Database(cfg.DBName))

	authH := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	noteH := handlers.NewNoteHandler(noteRepo, cfg)
	tagH := handlers.NewTagHandler(tagRepo)

	api := app.Group("/api")

	//helth
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Server running") })
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })

	// auth
	api.Post("/register", authH.Register)
	api.Post("/login", authH.Login)

	// notes (protected for create/update/delete)
	api.Post("/notes", middleware.RequireAuth(cfg), noteH.CreateNote)
	api.Get("/notes", middleware.RequireAuth(cfg), noteH.GetMyNotes)
	api.Get("/notes/public", noteH.GetPublicNotes)
	api.Get("/notes/:id", middleware.RequireAuth(cfg), noteH.GetNoteByID)
	api.Put("/notes/:id", middleware.RequireAuth(cfg), noteH.UpdateNote)
	api.Delete("/notes/:id", middleware.RequireAuth(cfg), noteH.DeleteNote)

	// tags
	api.Get("/tags/top", tagH.TopTags)

	return app
}
