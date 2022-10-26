package routes


import (
 
 
	//"github.com/labstack/echo/v4"
	"github.com/gofiber/fiber/v2"
	//"gorm.io/gorm"
	"cryptjoshi/handler"
	"gorm.io/gorm"
)



type Handler struct {
	db *gorm.DB
}


func  SetupRoute(e *fiber.App) {

	// user routes
 
	// e.Get("/",app.Root)
	e.Put("/user/:userId",handler.EditAUser)
	e.Post("/login",handler.UserLogin)
	e.Post("/admin",handler.AdminLogin)
	e.Get("/users", handler.GetAllUsers) 
	e.Get("/user/:userId",handler.GetAUser)
	e.Get("/user/balance/:userId", handler.UserBalance)
    e.Post("/register",handler.CreateUser)
	e.Delete("/user/:userId",handler.DeleteAUser)
	// e.Post("/user/passwd/:userId",handler.ChangePassword)
	// e.Get("/user/transaction/:userId",handler.Transaction)
	// e.Delete("/user/transaction/:userId",handler.ClearTransaction)
	e.Post("/user/deposit",handler.Deposit)
	// // provider routes
	e.Get("/operator/credit",handler.GetCreditOperator)
	// // game routes
	e.Get("/gamelist/:providerCode",handler.GetGamelist)
	e.Get("/gametype",handler.GetGametype)
	e.Get("/demogame/:provider/:gametype",handler.LuanchDGame)
	e.Get("/playgame/:userId/:provider/:gametype",handler.LuanchGame)
	// e.Put("/provider/:providerId",handler.EditAProvider)
	e.Get("/provider", handler.GetAllProviders) 
	// e.Get("/provider/:providerId",handler.GetAProvider)
	// e.Post("/addprovider",handler.CreateProvider)
	// e.Delete("/provider/:providerId",handler.DeleteAProivder)

	//seamless
	e.Post("/seamless/wallet/getuserbalance", handler.UserBalance)

}

