package utils

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/swagger/v2"
)

func UnmarshalRequest[T any](ctx fiber.Ctx) (T, error) {
	var payload T

	if err := json.Unmarshal(ctx.Body(), &payload); err != nil {
		return payload, err
	}

	return payload, nil
}

func InitDocs(router *fiber.App) {
	router.Get("/", func(ctx fiber.Ctx) error {
		return ctx.Status(fiber.StatusMovedPermanently).Redirect().To("/swagger/index.html")
	})

	router.Get("/swagger/*", swagger.HandlerDefault) // default

	router.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))
}
