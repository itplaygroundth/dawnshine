package handler

import (
	"net/http"
	"time"
	"cryptjoshi/configs"
	"cryptjoshi/models"
    "fmt"
   //"sort"
	//"cryptjoshi/responses"
	//"github.com/labstack/echo/v4"
    "cryptjoshi/utils"
    "strings"
	"cryptjoshi/responses"
    "github.com/tidwall/gjson"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"github.com/go-playground/validator/v10"
)


var providerCollection *mongo.Collection = configs.GetCollection(configs.MDB,"providers")
var validatep = validator.New()
 

func GetAProvider(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    providerId := c.Params("providerId")
    provider := new(models.Provider)
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(providerId)

    err := providerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(provider)

    if err != nil {
        return c.JSON(fiber.Map{"Status": http.StatusInternalServerError, "Message": "error", "Data": fiber.Map{"data": err.Error()}})
    }

    return c.JSON(fiber.Map{"Status": http.StatusOK, "Message": "success", "Data": fiber.Map{"data": provider}})
}

func CreateProvider(c *fiber.Ctx) error {	
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	provider := new(models.Provider)
	defer cancel()

	if err := c.BodyParser(provider); err!= nil {
		return c.JSON(fiber.Map{"Status": http.StatusBadRequest,"Message": "eror","Data": fiber.Map{"data":err.Error()}})
	}

	if validationErr := validatep.Struct(provider); validationErr!= nil {
		return c.JSON(fiber.Map{"Status": http.StatusBadRequest,"Message": "eror","Data": fiber.Map{"data":validationErr.Error()}})

	}

	newProvider := models.Provider{
		Id: primitive.NewObjectID(),
		ProviderCode:provider.ProviderCode,
		OperatorCode: configs.OperatorCode(),
		SecretKey:configs.SecretKey(),
		AgentCurrency:provider.AgentCurrency,
		BackendUrl:provider.BackendUrl,
		UserName:provider.UserName,
		Password:provider.Password,
	}

	result,err := providerCollection.InsertOne(ctx,newProvider)
	if err!= nil {
		return c.JSON(fiber.Map{"Status": http.StatusBadRequest,"Message": "eror","Data": fiber.Map{"data":err.Error()}})

	}
	return c.JSON(fiber.Map{"Status": http.StatusCreated,"Message": "success","Data": fiber.Map{"data":result}})		
}

func EditAProvider(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    providerId := c.Params("providerId")
    provider :=  new(models.Provider)
    defer cancel()

    objId, _ := primitive.ObjectIDFromHex(providerId)

    //validatep the request body
    if err := c.BodyParser(provider); err != nil {
        return c.JSON(fiber.Map{"Status": http.StatusBadRequest, "Message": "error", "Data": fiber.Map{"data": err.Error()}})
    }

    //use the validator library to validatep required fields
    if validationErr := validatep.Struct(provider); validationErr != nil {
        return c.JSON(fiber.Map{"Status": http.StatusBadRequest, "Message": "error", "Data": fiber.Map{"data": validationErr.Error()}})
    }

	update := bson.M{
		"providercode":provider.ProviderCode,
		"operatorcode":configs.OperatorCode(),
		"secretkey":configs.SecretKey(),
		"agentcurrency":provider.AgentCurrency,
		"backendurl":provider.BackendUrl,
		"username":provider.UserName,
		"password":provider.Password,
	}
	 

    result, err := providerCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

    if err != nil {
        return c.JSON(fiber.Map{"Status": http.StatusInternalServerError, "Message": "error", "Data": fiber.Map{"data": err.Error()}})
    }

    //get updated user details
    var updatedUser models.Provider
    if result.MatchedCount == 1 {
        err := providerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

        if err != nil {
            return c.JSON(fiber.Map{"Status": http.StatusInternalServerError, "Message": "error", "Data": fiber.Map{"data": err.Error()}})
        }
    }

    return c.JSON(fiber.Map{"Status": http.StatusOK, "Message": "success", "Data": fiber.Map{"data": updatedUser}})
}

func DeleteAProivder(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    providerId := c.Params("providerId")
    defer cancel()
  
    objId, _ := primitive.ObjectIDFromHex(providerId)
  
    result, err := providerCollection.DeleteOne(ctx, bson.M{"id": objId})
    if err != nil {
        return c.JSON(fiber.Map{"Status": http.StatusInternalServerError, "Message": "error", "Data": fiber.Map{"data": err.Error()}})
    }
  
    if result.DeletedCount < 1 {
        return c.JSON(fiber.Map{"Status": http.StatusNotFound, "Message": "error", "Data": fiber.Map{"data": "Provider with specified ID not found!"}})
    }
  
    return c.JSON(fiber.Map{"Status": http.StatusOK, "Message": "success", "Data": fiber.Map{"data": "Provider successfully deleted!"}})
}


func GetAllProviders(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var providers [] models.Provider
    defer cancel()

    results, err := providerCollection.Find(ctx, bson.M{})

    if err != nil {
        return c.JSON(fiber.Map{"Status": http.StatusInternalServerError, "Message": "error", "Data": fiber.Map{"data": err.Error()}})
    }

    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleProvider models.Provider
        if err = results.Decode(&singleProvider); err != nil {
            return c.JSON(fiber.Map{"Status": http.StatusInternalServerError, "Message": "error", "Data": fiber.Map{"data": err.Error()}})
        }

        providers = append(providers, singleProvider)
    }

    return c.JSON(fiber.Map{"Status": http.StatusOK, "Message": "success", "Data": fiber.Map{"data": providers}})
}

func GetCreditOperator(c *fiber.Ctx) error {
 
    signature := strings.ToUpper(string(utils.MD5(configs.OperatorCode()+configs.SecretKey())))

    url := configs.APIEndpoint()+"/checkAgentCredit.aspx?operatorcode="+configs.OperatorCode()+"&signature="+signature
 
    resp,_ := utils.FastGet(url,c)
   
   //fmt.Println(url)
   errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
   errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
   data := (gjson.Get(string(resp.Body()),"data")).String()

   if errCode != "0" {
      
       c.Status(fiber.StatusBadRequest)
       return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
    
   }

    return c.JSON(fiber.Map{"Status": http.StatusOK, "Message": "success", "Data": fiber.Map{"credit_balance": data}})
}

func GetGamelist(c *fiber.Ctx) error {
    
    provider := c.Params("providerCode")
    signature := strings.ToUpper(string(utils.MD5(configs.OperatorCode()+provider+configs.SecretKey())))

    url := configs.APIEndpoint()+"/getGameList.ashx?operatorcode="+configs.OperatorCode()+"&providercode="+provider+"&lang=en&html=0&reformatjson=yes&signature="+signature
 
    resp,_ := utils.FastGet(url,c)
   
   fmt.Println(url)
   errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
   errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
   gamelist := (gjson.Get(string(resp.Body()),"gamelist")).String()
  
   
   if errCode != "0" {
      
       c.Status(fiber.StatusBadRequest)
       return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
    
   }

    return c.JSON(fiber.Map{"Status": http.StatusOK, "Message": "success", "Data": fiber.Map{"gamelist": gamelist}})
}

func GetGametype(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    gametype := []models.Gametype{}
    defer cancel()

    result := configs.DB.WithContext(ctx).Find(&gametype)
    if result.Error != nil  {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"error": result.Error.Error()}})
    }

   return c.JSON(fiber.Map{"Status": http.StatusOK, "Message": "success", "Data": fiber.Map{"gametype": gametype}})
}

func LuanchGame(c *fiber.Ctx) error {
    //"WF, GT, IB"
    //provlist := []string{"WF","GT","IB"}
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    user := models.User{}
    gametype := c.Params("gametype")
    provider := c.Params("provider")
    userId := c.Params("userId")
   // gameid  := c.Params("gameId")
    // if ok := sort.SearchStrings(provlist,provider); ok != nil {
    //     return c.JSON(fiber.Map{"Status": http.StatusBadRequest, "Message": "error", "Data": fiber.Map{"gamelink": "Privider not available Game Demo"}})
    // }
    result := configs.DB.WithContext(ctx).First(&user,userId)
    if result.Error != nil  {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": result.Error.Error()}})
    }
    signature := strings.ToUpper(string(utils.MD5(configs.OperatorCode() + user.Password + strings.ToUpper(provider) + strings.ToUpper(gametype) + user.Name + configs.SecretKey())))
    url := configs.APIEndpoint()+"/launchGames.aspx?operatorcode="+configs.OperatorCode()+"&providercode="+strings.ToUpper(provider)+"&username="+user.Name+"&password="+user.Password+"&type="+strings.ToUpper(gametype)+"&signature="+signature
    //url := configs.APIEndpoint()+"/launchGames.aspx?operatorcode="+configs.OperatorCode()+"&providercode="+strings.ToUpper(provider)+"&username="+user.Name+"&password="+user.Password+"&type="+strings.ToUpper(gametype)+"&gameid="+gameid+"&signature="+signature
 
    resp,_ := utils.FastGet(url,c)
   
   //fmt.Println(url)
   errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
   errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
   data := (gjson.Get(string(resp.Body()),"gameUrl")).String()

   if errCode != "0" {
      
       c.Status(fiber.StatusBadRequest)
       return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
    
   }

    return c.JSON(fiber.Map{"Status": http.StatusOK, "Message": "success", "Data": fiber.Map{"gameUrl": data}})

}
func LuanchDGame(c *fiber.Ctx) error {
    //"WF, GT, IB"
    provlist := []string{"WF","GT","IB"}
   
    gametype := strings.ToUpper(c.Params("gametype"))
    provider := strings.ToUpper(c.Params("provider"))
  
    if ok := utils.Contains(provlist,provider); ok == false {
     fmt.Println(ok)
        return c.JSON(fiber.Map{"Status": http.StatusBadRequest, "Message": "error", "Data": fiber.Map{"gamelink": "Provider not available Game Demo"}})
    }
    url := configs.APIEndpoint()+"/launchDGames.ashx?providercode="+provider+"&type="+gametype
 
    resp,_ := utils.FastGet(url,c)
    //fmt.Println(url)
   
   errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
   errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
   gameUrl := (gjson.Get(string(resp.Body()),"gameUrl")).String()

   if errCode != "0" {
      
       c.Status(fiber.StatusBadRequest)
       return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
    
   }

    return c.JSON(fiber.Map{"Status": http.StatusOK, "Message": "success", "Data": fiber.Map{"gamelink": gameUrl}})

}
