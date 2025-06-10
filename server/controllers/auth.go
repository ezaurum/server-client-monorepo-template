package controllers

import (
	"github.com/gofiber/fiber/v2"
	"templateapp/services"
)

type AuthController struct {
	Service *services.UserService
}

func (c AuthController) Login(ctx *fiber.Ctx) error {
	return nil
}
