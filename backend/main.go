package main

import (
    //"github.com/labstack/echo/v4"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"cryptjoshi/configs"
	"cryptjoshi/routes"
)



func main() {

	app := fiber.New()
	app.Use(favicon.New())
	app.Use(cors.New())
	
	app.Use(cors.New(cors.Config{
	  AllowOrigins: "*",
	  AllowHeaders:  "Origin, Content-Type, Accept",
	  AllowMethods:"GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	configs.ConnectDB()
	app.Use(logger.New())

	routes.SetupRoute(app)
    log.Fatal(app.Listen(":3333"))
}