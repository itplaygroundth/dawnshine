package handler

import (
	"net/http"
    "fmt"
	 "time"
    // "encoding/json"
	"cryptjoshi/configs"
	"cryptjoshi/models"
    "cryptjoshi/utils"
    "strings"
	"cryptjoshi/responses"
    "github.com/gofiber/fiber/v2"
	//"github.com/labstack/echo/v4"
    "github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"github.com/go-playground/validator/v10"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB,"users")
var validate = validator.New()

func Root(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).SendString("Hello there!")
}

func GetAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    user := new(models.User)
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(userId)

    err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

    if err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}} )
             }

    return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}



func CreateUser(c *fiber.Ctx) error {	
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
 
    user := new(models.User)
	defer cancel()

	if err := c.BodyParser(user); err!= nil {
	       return c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "eror",Data: &fiber.Map{"data":err.Error()}})
	
    }

	if validationErr := validate.Struct(user); validationErr!= nil {
		return c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "eror",Data: &fiber.Map{"data":validationErr.Error()}})

	}
    user.Signature = strings.ToUpper(string(utils.MD5(configs.OperatorCode()+user.Name+configs.SecretKey())))
  
	newUser := models.User{
		Id: primitive.NewObjectID(),
		Name: user.Name,
        Location: user.Location,
		Title: user.Title,
        Password: user.Password,
        Balance: 0.0,
        OperatorCode: configs.OperatorCode(),
        Signature: user.Signature,
	}

    url := configs.APIEndpoint()+"/createMember.aspx?operatorcode="+configs.OperatorCode()+"&username="+user.Name+"&signature="+user.Signature
    
    resp,_ := utils.FastGet(url,c)
    
 
    errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
    errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
    if errCode != "0" {
        c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
     
    }

	result,err := userCollection.InsertOne(ctx,newUser)
	if err!= nil {
		c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "eror",Data: &fiber.Map{"data":err.Error()}}))

	}
  
    
    return c.JSON(responses.UserResponse{Status: http.StatusCreated,Message: "success",Data: &fiber.Map{"data":result}})		

}

func EditAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    user := new(models.User)
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(userId)

    //validate the request body
    if err := c.BodyParser(user); err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(user); validationErr != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    }
    
     

    update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title,"password":user.Password,"operatorcode": configs.OperatorCode(),
    "signature": utils.MD5(configs.OperatorCode()+user.Name+configs.SecretKey()),}

    result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

    if err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //get updated user details
    var updatedUser models.User
    
    if result.MatchedCount == 1 {
        err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

        if err != nil {
            return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
        }
    }

    return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    defer cancel()
  
    objId, _ := primitive.ObjectIDFromHex(userId)
  
    result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
    if err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }
  
    if result.DeletedCount < 1 {
        return c.JSON(responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}})
    }
  
    return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}})
}


func GetAllUsers(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var users []models.User
    defer cancel()

 

    results, err := userCollection.Find(ctx, bson.M{})

    if err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleUser models.User
        if err = results.Decode(&singleUser); err != nil {
            return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
        }

        users = append(users, singleUser)
    }

    return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}})
}


func ChangePassword(c *fiber.Ctx) error{
 
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    user := new(models.User)
    nuser := new(models.User)
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(userId)

    err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

    if err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})

    }

    opassword := user.Password
    username := user.Name
  
    if err := c.BodyParser(nuser); err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }

    //use the validator library to validate required fields
    // if validationErr := validate.Struct(nuser); validationErr != nil {
    //     return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
    // }
   
    //  response := checkUserProduct(username,c)
    //  jsonStr, _ := json.Marshal(response.Data)
    //  fmt.Println(gjson.Get(string(jsonStr),"Data"))
     

    provider := strings.ToUpper(c.Subdomains(2)[0])
    signature := strings.ToUpper(string(utils.MD5(opassword + configs.OperatorCode() + nuser.Password + strings.ToUpper(c.Subdomains(2)[0]) + username + configs.SecretKey())))
    url := configs.APIEndpoint()+"/changePassword.aspx?operatorcode="+configs.OperatorCode()+"&providercode="+provider+"&username="+username+"&password="+nuser.Password+"&opassword="+opassword+"&signature="+signature
    fmt.Println(url)
    resp,_ := utils.FastGet(url,c)
    errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
    errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
    if errCode != "0" {
        c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
     
    }
    return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}

func UserBalance(c *fiber.Ctx) error {
 
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    user := new(models.User)
    defer cancel()

  

    objId, _ := primitive.ObjectIDFromHex(userId)

    err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

    if err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})

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

func checkUserProduct(username string,c *fiber.Ctx) responses.ProductResponse {
    _, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    // userId := c.Params("userId")
    // user := new(models.User)
    defer cancel()

  

    // objId, _ := primitive.ObjectIDFromHex(userId)

    // err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

    // if err != nil {
    //     return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})

    // }
    provider := strings.ToUpper(c.Subdomains(2)[0])
    hashstr := configs.OperatorCode()+provider+username+configs.SecretKey()

    signature := strings.ToUpper(string(utils.MD5(hashstr)))
    url := configs.APIEndpoint()+"/checkMemberProductUsername.ashx?operatorcode="+configs.OperatorCode()+"&providercode="+provider+"&username="+username+"&signature="+signature
    
    resp,_ := utils.FastGet(url,c)
     
    

    errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
    errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()

    if errCode != "0" {
        //c.Status(fiber.StatusBadRequest)
        return responses.ProductResponse{Status: http.StatusBadRequest,Message: errMsg,Data: fiber.Map{"errCode": errCode,"errMsg":errMsg}}
     
    }
       return responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: fiber.Map{"Data": (gjson.Get(string(resp.Body()),"data")).String()}}
   
    }