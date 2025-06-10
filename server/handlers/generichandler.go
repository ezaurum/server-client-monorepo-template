package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"templateapp/database"
)

type JsonListResponse[T any] struct {
	List []T `json:"list"`
}

func List[T any](c *fiber.Ctx) error {
	var list []T
	if result := database.DB().Find(&list); result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	return c.JSON(
		JsonListResponse[T]{
			List: list,
		})
}

func Read[T any](c *fiber.Ctx) error {
	var item T
	if result := database.DB().First(&item, c.Params("id")); result.Error != nil {
		return result.Error
	}
	return c.JSON(item)
}

func Create[T any](c *fiber.Ctx) error {
	var item T
	if err := c.BodyParser(&item); err != nil {
		return err
	}
	if result := database.DB().Create(&item); result.Error != nil {
		return result.Error
	}
	return c.JSON(item)
}

func Update[T any](c *fiber.Ctx) error {
	var item T
	if err := c.BodyParser(&item); err != nil {
		return err
	}
	if result := database.DB().Save(&item); result.Error != nil {
		return result.Error
	}
	return c.JSON(item)
}

func Delete[T any](c *fiber.Ctx) error {
	var item T
	if result := database.DB().Delete(&item, c.Params("id")); result.Error != nil {
		return result.Error
	}
	return c.SendStatus(fiber.StatusNoContent)
}
