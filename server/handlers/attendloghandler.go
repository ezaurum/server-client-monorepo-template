package handlers

import (
	"github.com/gofiber/fiber/v2"
	"templateapp/database"
	"templateapp/models"
)

func ListCheckInLog(c *fiber.Ctx) error {
	c.Query("gatheringID")
	database.DB().Find(&models.CheckInLog{})
	return c.SendString("ListAttendLog")
}
