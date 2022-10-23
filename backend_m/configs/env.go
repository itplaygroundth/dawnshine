package configs

import (
	// "log"
	"os"
	// "github.com/joho/godotenv"
)

func EnvMongoURI() string {
	// err := godotenv.Load()
	// if err!= nil {
    //     log.Fatal(err)
    // }
	return os.Getenv("MONGODB")
}

 
func DbConnection() string {
	return os.Getenv("DB_CONNECTION")
}

func OperatorCode() string {
	return os.Getenv("OPERATORCODE")
}
func APIEndpoint() string {
	return os.Getenv("API_ENDPOINT")
}

func LOGEndpoint() string {
	return os.Getenv("LOG_ENDPOINT")
}

func SecretKey() string {	
	return os.Getenv("SECRET_KEY")
}