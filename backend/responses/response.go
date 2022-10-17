package responses

import (
    //"github.com/labstack/echo/v4"
    "github.com/gofiber/fiber/v2"
)
type UserResponse struct {
    Status  int       `json:"status"`
    Message string    `json:"message"`
    Data    *fiber.Map `json:"data"`
}

type ProviderResponse struct {
    Status  int       `json:"status"`
    Message string    `json:"message"`
    Data    *fiber.Map `json:"data"`
}

type ProductResponse struct {
    Status  int       `json:"status"`
    Message string    `json:"message"`
    Data    fiber.Map `json:"data"`
}
