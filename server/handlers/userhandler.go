package handlers

import (
	"encoding/gob"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"templateapp/database"
	"templateapp/models"
)

// UserList returns a list of users
func UserList(c *fiber.Ctx) error {
	users := database.Get()

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users,
	})
}

// UserCreate registers a user
func UserCreate(c *fiber.Ctx) error {
	user := &models.User{
		// Note: when writing to external database,
		// we can simply use - Name: c.FormValue("user")
		Name: utils.CopyString(c.FormValue("user")),
	}
	database.Insert(user)

	log.Println("User created:", user.Name)

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

// UserGet returns a user
// @Summary Get a user
// @Description Get a user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func UserGet(ctx *fiber.Ctx) error {
	var data = make(map[string]interface{})
	data["id"] = ctx.Params("id")
	data["name"] = "John Doe"
	return ctx.JSON(data)
}

func UserDelete(c *fiber.Ctx) error {
	return c.SendString("User deleted")
}

func init() {
	// User 타입 등록
	gob.Register(models.User{})
}
