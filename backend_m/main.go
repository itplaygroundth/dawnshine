package main

import (
    //"github.com/labstack/echo/v4"
	//"net/http"
	"os"
	"gorm.io/driver/mysql"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/swaggo/swag/cmd/swag"
    // "github.com/arsmn/fiber-swagger/v2"
	"cryptjoshi/configs"
    "cryptjoshi/users"
	"cryptjoshi/routes" 
 	"gorm.io/gorm"
	//"cryptjoshi/handler"

)
// @title APIs
// @version 1.0
// @description This is a sample server server.
// @termsOfService https://bkkdragon.magic/

// @contact.name API Support
// @contact.url https://bkkdragon.magic/
// @contact.email drfate@bkkdragon.magic

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host code.lab.lan
// @BasePath /
// @schemes https

func initDB() *gorm.DB{
	dial := mysql.Open(os.Getenv("DB_CONNECTION"))
	db, err := gorm.Open(dial)
	if err != nil {
		panic(err)
	}
	 

	db.AutoMigrate(&users.User{})

	return db
}

func main() {
	//configs.ConnectDB()
	//db := initDB()
	//defer db.Close()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(favicon.New())
 
 
	app.Use(cors.New(cors.Config{
	  AllowOrigins: "*",
	  AllowHeaders:  "Origin, Content-Type, Accept",
	  AllowMethods:"GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	 
	configs.DB.AutoMigrate(users.User{})
	configs.DB.AutoMigrate(users.Balance{})
	configs.DB.AutoMigrate(users.Wallet{})
	configs.DB.AutoMigrate(users.Gametype{})
	
	//userAPI := initUserAPI(configs.DB)
 
	//app.Get("/users",userAPI.FindAll)
	// app.Get("/users/:userid",UserAPI.FindByID)
	// app.Post("/users",UserAPI.CreateUser)
	// app.Put("/users/:userid",UserAPI.Edit)
	// app.Delete("/users/:userid",UserAPI.Delete)
	//log.Fatal(&handler)
	// app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
	// 	URL:         "/swagger/doc.json",
	// 	DeepLinking: false,
	// }))
    routes.SetupRoute(app)
 
    log.Fatal(app.Listen(":3333"))
}
 