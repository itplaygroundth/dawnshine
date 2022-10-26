package models

//import "go.mongodb.org/mongo-driver/bson/primitive"
import "gorm.io/gorm"
//import "gorm.io/datatypes"

type User struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement:true"`
	Name string  `gorm:"primaryKey"`    
	Location string  
	Title	 string 
	Password string  
	Balance float32  
	OperatorCode string  
	Signature string
	Isadmin bool `gorm:"type:bool;default:false"`
}

type Balance struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement:true"`
	Transactionid	string	 `json:"transactionid"`	
	Opcode	string	 `json:"opcode"`	
	Userid	uint64	 `json:"userid"`	
	Amount	float32	 `json:amount"`
	Currency string	`json:currency"`
	Password	string	`json:password"`		
}

type Wallet struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement:true"`
	Method string  
	OperatorCode string  
	Providercode	 string 
	Userid uint64  
	Username string  
	Password string  
	Referenceid string  
	Action	 string  
	Amount     float32  	
	Signature string 
	Time string  
	Status string 
}

type Gametype struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement:true"`
	Code string `json:"code,omitempty" validate:"required"`
	Name string `json:"name,omitempty" validate:"required"`
}