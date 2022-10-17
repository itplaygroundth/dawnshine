package handler

import (
	"net/http"
	"time"
	"cryptjoshi/configs"
	"cryptjoshi/models"
	//"cryptjoshi/responses"
	//"github.com/labstack/echo/v4"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"github.com/go-playground/validator/v10"
)


var providerCollection *mongo.Collection = configs.GetCollection(configs.DB,"providers")
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