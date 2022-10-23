package configs

import (
	"context"
	"fmt"
	"log"
	"time"
	 
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MConnectDB() *mongo.Client {
	client,err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err!= nil {
        log.Fatal(err)
    }

	ctx,_ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
    if err!= nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx,nil)
	if err!= nil {
        log.Fatal(err)
    }
	fmt.Println("Connect to Mongo successfully")
	return client
}


func ConnectDB() *gorm.DB {

	//dsn := "root:12345678@tcp(localhost:3306)/go_basics?parseTime=true"
	dial := mysql.Open(DbConnection())
	db, err := gorm.Open(dial)
	if err != nil {
		panic(err)
	}

	// db,err := sql.Open("mysql",DbConnection())
	 
	// if err!= nil {
    //     log.Fatal(err)
    // }
	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)
	//ctx,_ := context.WithTimeout(context.Background(), 10*time.Second)
	// err = client.Connect(ctx)
    // if err!= nil {
	// 	log.Fatal(err)
	// }

	// err = client.Ping(ctx,nil)
	// if err!= nil {
    //     log.Fatal(err)
    // }
	fmt.Println("Connect to Mysql successfully")
	return db
}
var MDB *mongo.Client = MConnectDB()
var DB *gorm.DB = ConnectDB()

func GetCollection(client *mongo.Client,collectionName string) *mongo.Collection {
	collection := client.Database("sunshinebet").Collection(collectionName)
	return collection
}

