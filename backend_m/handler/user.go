package handler

import (
	"net/http"
 
    "time"
    "strconv"
    
 
	"cryptjoshi/configs"
	"cryptjoshi/users"
    "cryptjoshi/utils"
    "strings"
	"cryptjoshi/responses"
    "github.com/gofiber/fiber/v2"
 
    "github.com/tidwall/gjson"
    "gorm.io/gorm"
	"golang.org/x/net/context"
	"github.com/go-playground/validator/v10"
)
 
var validate = validator.New()
 
type Handler struct {
	db *gorm.DB
}

 

func  Root(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).SendString("Hello there!")
}

 func  GetAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    userId := c.Params("userId")    
 
    user := users.User{}
 
   
    result := configs.DB.WithContext(ctx).First(&user,userId)
    if result.Error != nil  {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": result.Error.Error()}})
    }
  
     return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
 }

 func   GetAllUsers(c *fiber.Ctx) error {
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
 
    defer cancel()
    user := []users.User{}
 
    
    result := configs.DB.WithContext(ctx).Find(&user)
    if result.Error != nil  {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": result.Error.Error()}})
    }
 
     return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
 }



 func   CreateUser(c *fiber.Ctx) error {	
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
   defer cancel()
   user := users.User{}
 
  

	if err := c.BodyParser(&user); err!= nil {
	       return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError,Message: "eror",Data: &fiber.Map{"data":err.Error()}})
	
    }
    if validationErr := validate.Struct(user); validationErr!= nil {
		return c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "eror",Data: &fiber.Map{"data":validationErr.Error()}})

	}
    user.Signature = strings.ToUpper(string(utils.MD5(configs.OperatorCode()+user.Name+configs.SecretKey())))
 	newUser := users.User{
		Name: user.Name,
        Location: user.Location,
		Title: user.Title,
        Password: user.Password,
        Balance: 0.0,
        OperatorCode: configs.OperatorCode(),
        Signature: user.Signature,
	}
    tx := configs.DB.WithContext(ctx).Begin()
    defer func() {
        if r := recover(); r != nil {
          tx.Rollback()
        }
      }()
    
      if err := tx.Error; err != nil {
        c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "error",Data: &fiber.Map{"errCode": err,"errMsg":err}}))

      }
    
    if  err := tx.Create(&newUser).Error; err != nil {
        tx.Rollback()
        c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "error",Data: &fiber.Map{"errCode": err,"errMsg":err}}))

    }
  
    

     url := configs.APIEndpoint()+"/createMember.aspx?operatorcode="+configs.OperatorCode()+"&username="+user.Name+"&signature="+user.Signature
 
     resp,_ := utils.FastGet(url,c)
    
 
    errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
    errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
    if errCode != "0" {
        tx.Rollback()
        c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
     
    }
    tx.Commit()
    return c.JSON(responses.UserResponse{Status: http.StatusCreated,Message: "success",Data: &fiber.Map{"data":newUser}})		

 }

func   EditAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    userId,_ := strconv.Atoi(c.Params("userId"))    
    user := users.User{}
    result := configs.DB.WithContext(ctx).First(&user,userId)
    if result.Error != nil  {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": result.Error.Error()}})
    }
    //validate the request body
    if err := c.BodyParser(&user); err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(user); validationErr != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    }
    
 
   
    updateUser := users.User{
        ID: uint64(userId),
		//Name: user.Name,
        Location: user.Location,
		Title: user.Title,
        Password: user.Password,
        Balance: 0.0,
       // OperatorCode: configs.OperatorCode(),
       // Signature: utils.MD5(configs.OperatorCode()+user.Name+configs.SecretKey()),
	}
      result = configs.DB.WithContext(ctx).Updates(updateUser)

    
    if result.Error != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": result.Error}})
    }
    
    //get updated user details
    return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updateUser}})
}

func   DeleteAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    defer cancel()
  
    user := users.User{}
    
   
    result := configs.DB.WithContext(ctx).Delete(&user,userId)
 
    if result.Error != nil  {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": result.Error.Error()}})
    }
  
 
    
  
    return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}})
}

 


// func ChangePassword(c *fiber.Ctx) error{
 
//     ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//     userId := c.Params("userId")
//     user := new(users.User)
//     nuser := new(users.User)
//     defer cancel()

//     objId, _ := primitive.ObjectIDFromHex(userId)

//     err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

//     if err != nil {
//         return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})

//     }

//     opassword := user.Password
//     username := user.Name
  
//     if err := c.BodyParser(nuser); err != nil {
//         return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
//     }

//     //use the validator library to validate required fields
//     // if validationErr := validate.Struct(nuser); validationErr != nil {
//     //     return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
//     // }
   
//     //  response := checkUserProduct(username,c)
//     //  jsonStr, _ := json.Marshal(response.Data)
//     //  fmt.Println(gjson.Get(string(jsonStr),"Data"))
     

//     provider := strings.ToUpper(c.Subdomains(2)[0])
//     signature := strings.ToUpper(string(utils.MD5(opassword + configs.OperatorCode() + nuser.Password + strings.ToUpper(c.Subdomains(2)[0]) + username + configs.SecretKey())))
//     url := configs.APIEndpoint()+"/changePassword.aspx?operatorcode="+configs.OperatorCode()+"&providercode="+provider+"&username="+username+"&password="+nuser.Password+"&opassword="+opassword+"&signature="+signature
//     fmt.Println(url)
//     resp,_ := utils.FastGet(url,c)
//     errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
//     errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
//     if errCode != "0" {
//         c.Status(fiber.StatusBadRequest)
//         return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
     
//     }
//     return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
// }

func UserBalance(c *fiber.Ctx) error {
 
 

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    userId := c.Params("userId")    
 
    user := users.User{}
 
   
    result := configs.DB.WithContext(ctx).First(&user,userId)
    if result.Error != nil  {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": result.Error.Error()}})
    }

    provider := strings.ToUpper(c.Subdomains(2)[0])
    hashstr := configs.OperatorCode()+user.Password+provider+user.Name+configs.SecretKey()

    signature := strings.ToUpper(string(utils.MD5(hashstr)))
    url := configs.APIEndpoint()+"/getBalance.aspx?operatorcode="+configs.OperatorCode()+"&providercode=AG&username="+user.Name+"&password="+user.Password+"&opassword=&signature="+signature
    
    resp,_ := utils.FastGet(url,c)
 
    errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
    errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
 
    if errCode != "0" {
        c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
     
    }
    return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"balance": (gjson.Get(string(resp.Body()),"balance")).String()}})
}

func CheckUserProduct(username string,c *fiber.Ctx) responses.ProductResponse {
 
    
    provider := strings.ToUpper(c.Subdomains(2)[0])
    hashstr := configs.OperatorCode()+provider+username+configs.SecretKey()

    signature := strings.ToUpper(string(utils.MD5(hashstr)))
    url := configs.APIEndpoint()+"/checkMemberProductUsername.ashx?operatorcode="+configs.OperatorCode()+"&providercode="+provider+"&username="+username+"&signature="+signature
    
    resp,_ := utils.FastGet(url,c)
     
    

    errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
    errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()

    if errCode != "0" {
 
        return responses.ProductResponse{Status: http.StatusBadRequest,Message: errMsg,Data: fiber.Map{"errCode": errCode,"errMsg":errMsg}}
     
    }
       return responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: fiber.Map{"Data": (gjson.Get(string(resp.Body()),"data")).String()}}
   
    }