package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"templateapp/database"
	"templateapp/models"
)

// GenericList returns a list of users
func GenericList(c *fiber.Ctx) error {
	users := database.Get()

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users,
	})
}

// GenericCreate registers a user
func GenericCreate(c *fiber.Ctx) error {
	user := &models.User{
		// Note: when writing to external database,
		// we can simply use - Name: c.FormValue("user")
		Name: utils.CopyString(c.FormValue("user")),
	}
	database.Insert(user)

	log.Println("Generic created:", user.Name)

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}

// GenericGet returns a user
// @Summary Get a user
// @Description Get a user
// @Tags users
// @Produce json
// @Param id path int true "Generic ID"
// @Success 200 {object} models.Generic
// @Router /users/{id} [get]
func Get(ctx *fiber.Ctx) error {
	var data = make(map[string]interface{})
	data["id"] = ctx.Params("id")
	data["name"] = "John Doe"
	return ctx.JSON(data)
}

func GenericDelete(c *fiber.Ctx) error {
	return c.SendString("Generic deleted")
}
