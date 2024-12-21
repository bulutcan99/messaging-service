package handlers

import "github.com/gofiber/fiber/v3"

func GetRoute1(c fiber.Ctx) error {
	return c.SendString("This is GET route 1")
}

func GetRoute2(c fiber.Ctx) error {
	return c.SendString("This is GET route 2")
}

func PostRoute1(c fiber.Ctx) error {
	return c.SendString("This is POST route 1")
}

func PostRoute2(c fiber.Ctx) error {
	return c.SendString("This is POST route 2")
}
