package handler

import (
	"net/http"
 
    "fmt"
	"time"
 
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
	// "github.com/go-playground/validator/v10"
)

var walletCollection *mongo.Collection = configs.GetCollection(configs.DB,"wallet")
//var validate = validator.New()

type Balance struct {
	Code	string	 `json:"code"`	
	Opcode	string	 `json:"opcode"`	
	Userid	string	 `json:"userid"`	
	Balance	float32	 `json:balance"`
	Password	string	`json:password"`		
}

func GetUserBalance(c *fiber.Ctx) error {
 
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
		url := configs.APIEndpoint()+"/getBalance.aspx?operatorcode="+configs.OperatorCode()+"&providercode="+provider+"&username="+user.Name+"&password="+user.Password+"&signature="+signature
		
		resp,_ := utils.FastGet(url,c)
		fmt.Println(string(resp.Body()))
		errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
		errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
		if errCode != "0" {
			c.Status(fiber.StatusBadRequest)
			return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
		 
		}
		mod := mongo.IndexModel{
			Keys: bson.M{
			"userid": 1, // index in ascending order
			}, Options: nil,
			}
		ind, err = col.Indexes().CreateOne(ctx, mod)
		if err != nil {
			fmt.Println("Indexes().CreateOne() ERROR:", err)
			os.Exit(1) // exit in case of error
			} else {
			// API call returns string of the index name
			fmt.Println("CreateOne() index:", ind)
			fmt.Println("CreateOne() type:", reflect.TypeOf(ind), "\n")
			}
		return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": resp.Body()}})
}

func UpdateTransaction(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)
	update:=  bson.M{"status": "complete"}
	result,err := walletCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	 
	if err!= nil {
		c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "eror",Data: &fiber.Map{"data":err.Error()}}))

	}
	 
    return c.JSON(responses.UserResponse{Status: http.StatusCreated,Message: "success",Data: &fiber.Map{"data":result}})
    
}

func ClearTransaction(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    defer cancel()

 
	result,err := walletCollection.DeleteMany(ctx, bson.M{"userid": userId})
	if err!= nil {
		c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "eror",Data: &fiber.Map{"data":err.Error()}}))

	}
	 
    return c.JSON(responses.UserResponse{Status: http.StatusCreated,Message: "success",Data: &fiber.Map{"data":result}})
    
}
func Transaction(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    defer cancel()


    cursor,err := walletCollection.Find(ctx, bson.M{"userid": userId})
	var wallet []bson.M
	if err = cursor.All(ctx, &wallet); err != nil {
		  
         return (c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "eror",Data: &fiber.Map{"data":err.Error()}}))

	}
    if err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}} )
             }

    return c.JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": wallet}})
}

func Deposit(c *fiber.Ctx) error {	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	wallet := new(models.Wallet)
	user := new(models.User)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

	if err != nil {
		return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})

	}

	password := user.Password
    username := user.Name
	provider := strings.ToUpper(c.Subdomains(2)[0])

    if err := c.BodyParser(wallet); err != nil {
        return c.JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
    }
	hashstr := fmt.Sprintf("%f",wallet.Amount)+password+provider+wallet.Referenceid+"deposit"+username+configs.SecretKey()
	signature := strings.ToUpper(string(utils.MD5(hashstr)))
	now := time.Now()
	newDeposit := models.Wallet{
		Id: primitive.NewObjectID(),
		Method:"GET",
		OperatorCode: configs.OperatorCode(),
		Providercode:provider,
		Userid: userId,
		Username: username,
		Password: password,
		Referenceid: wallet.Referenceid,
		Action: "deposit",
		Amount: wallet.Amount,
		Signature: signature,
		Time: fmt.Sprintf("%d",now.Unix()),
		Status: "wait",
	}
	//amount + operatorcode + password + providercode + referenceid + type + username + secret_key
	_,err = walletCollection.InsertOne(ctx,newDeposit)
	if err!= nil {
		c.Status(fiber.StatusBadRequest)
        return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: "eror",Data: &fiber.Map{"data":err.Error()}}))

	}
	url := configs.APIEndpoint()+"/makeTransfer.aspx?operatorcode="+configs.OperatorCode()+"&providercode="+provider+"&username="+username+"&password="+password+"&referenceid="+wallet.Referenceid+"&type=deposit&amount="+fmt.Sprintf("%f",wallet.Amount)+"&signature="+signature
	fmt.Println(url)
	resp,_ := utils.FastGet(url,c)
	// fmt.Println(string(resp.Body()))
	errCode := (gjson.Get(string(resp.Body()),"errCode")).String()
	errMsg  := (gjson.Get(string(resp.Body()),"errMsg")).String()
	 if errCode != "0" {
		c.Status(fiber.StatusBadRequest)
		return (c.JSON(responses.UserResponse{Status: http.StatusBadRequest,Message: errMsg,Data: &fiber.Map{"errCode": errCode,"errMsg":errMsg}}))
	 
	}
	
    return c.JSON(responses.UserResponse{Status: http.StatusCreated,Message: "success",Data: &fiber.Map{"data":resp}})
}

func Withdraw(c *fiber.Ctx) {    

}

