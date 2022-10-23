package routes


import (
	"cryptjoshi/handler"
	//"github.com/labstack/echo/v4"
	"github.com/gofiber/fiber/v2"
	
)

func SetupRoute(e *fiber.App) {
	// user routes
	e.Get("/",handler.Root)
	e.Put("/user/:userId",handler.EditAUser)
	e.Get("/users", handler.GetAllUsers) 
	e.Get("/user/:userId",handler.GetAUser)
	e.Get("/user/balance/:userId", handler.UserBalance)
	e.Post("/adduser",handler.CreateUser)
	e.Delete("/user/:userId",handler.DeleteAUser)
	e.Post("/user/passwd/:userId",handler.ChangePassword)
	e.Get("/user/transaction/:userId",handler.Transaction)
	e.Delete("/user/transaction/:userId",handler.ClearTransaction)
	e.Post("/user/deposit/:userId",handler.Deposit)
	// provider routes
	e.Put("/provider/:providerId",handler.EditAProvider)
	e.Get("/providers", handler.GetAllProviders) 
	e.Get("/provider/:providerId",handler.GetAProvider)
	e.Post("/addprovider",handler.CreateProvider)
	e.Delete("/provider/:providerId",handler.DeleteAProivder)

	//seamless
	e.Post("/seamless/wallet/getuserbalance", handler.GetUserBalance)

}

