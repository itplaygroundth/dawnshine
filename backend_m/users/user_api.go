package users

import (
	"net/http"
	"strconv"
    //"fmt"
    // "log"
    //"time"
    //"encoding/json"
	"cryptjoshi/configs"
 
    "cryptjoshi/utils"
    "strings"
	"cryptjoshi/responses"
    "github.com/gofiber/fiber/v2"
	//"github.com/labstack/echo/v4"
    "github.com/tidwall/gjson"
    //"gorm.io/gorm"
	//"golang.org/x/net/context"
	//"github.com/go-playground/validator/v10"
)

type UserAPI struct {
	UserService UserService
}

func ProvideUserAPI(u UserService) UserAPI {
	return UserAPI{UserService: u}
}

func (u *UserAPI) FindAll(c *fiber.Ctx) error {
	
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    //userId := c.Params("userId")    
    // user := new(User)
    //defer cancel()
	users := u.UserService.FindAll()
	// if users == (User{})  {
    //     return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": users.Error.Error()}})
    // }
 
     return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}})

}

func (u *UserAPI) FindByID(c *fiber.Ctx) error {
	id,_ := strconv.Atoi(c.Params("userid"))
	users := u.UserService.FindByID(uint64(id))
	 
     return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}})

}

func (u *UserAPI)  CreateUser(c *fiber.Ctx) error {	

	user :=  User{}
	
	 if err := c.BodyParser(&user); err!= nil {
			return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError,Message: "eror",Data: &fiber.Map{"data":err.Error()}})
	 
	 }
	
	 
	 user.Signature = strings.ToUpper(string(utils.MD5(configs.OperatorCode()+user.Name+configs.SecretKey())))
	 newUser := User{
		 Name: user.Name,
		 Location: user.Location,
		 Title: user.Title,
		 Password: user.Password,
		 Balance: 0.0,
		 OperatorCode: configs.OperatorCode(),
		 Signature: user.Signature,
	 }
	 
	 createUser := u.UserService.Save(newUser)
	 
 
	  url := configs.APIEndpoint()+"/createMember.aspx?operatorcode="+configs.OperatorCode()+"&username="+user.Name+"&signature="+user.Signature
  
	  resp,_ := utils.FastGet(url,c)
	 
  
	 errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
	 errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
	 if errCode != "0" {
	 
		 c.Status(fiber.StatusBadRequest)
		 return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
	  
	 }
 
	 return c.JSON(responses.UserResponse{Status: http.StatusCreated,Message: "success",Data: &fiber.Map{"data": createUser}})		
 
  }
 
  func (u *UserAPI) Edit(c *fiber.Ctx) error {
	 
	userId,_ := strconv.Atoi(c.Params("userId"))
	user := u.UserService.FindByID(uint64(userId))

	 if user == (User{}) {
		return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": user}})
	 }

	 userDto :=  User{}
	 //validate the request body
	 if err := c.BodyParser(&userDto); err != nil {
		 return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	 }
 
	 //use the validator library to validate required fields
	  
	 updateUser := User{
		 Name: userDto.Name,
		 Location: userDto.Location,
		 Title: userDto.Title,
		 Password: userDto.Password,
		 Balance: 0.0,
		 OperatorCode: configs.OperatorCode(),
		 Signature: utils.MD5(configs.OperatorCode()+userDto.Name+configs.SecretKey()),
	 }
	 
	  u.UserService.Save(updateUser)
	 
	 
	 
	 //get updated user details
	 return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updateUser}})
 }
 
 func (u *UserAPI)  Delete(c *fiber.Ctx) error {
	  userId,_ := strconv.Atoi(c.Params("userId"))
	  user := u.UserService.FindByID(uint64(userId))

	  if user == (User{}){
		return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": user}})

	  }
	  u.UserService.Delete(user)
	 
	 return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}})
 }