package routes

import (
	"net/http"
	//"log"
	"github.com/gofiber/fiber/v2"

	// "github.com/gofiber/fiber/v2/middleware/favicon"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	//"cryptjoshi/configs"
	//"cryptjoshi/routes" 
	//"cryptjoshi/utils"
 
	"gorm.io/gorm"
)



type Handler struct {
	db *gorm.DB
}
// type Controller struct {
// 	service *utils.Service 
// }
// type ControlInterface struct {
// 	c *fiber.App
// 	service *utils.Service
// }

// func NewController(control ControlInterface) *Controller {
//     return &Controller{control}
// }

func (h *Handler)  Root(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).SendString("Hello there!")
}