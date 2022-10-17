package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id primitive.ObjectID `json:"id,omitempty"`
	Name string            `json:"name,omitempty" validate:"required"`
	Location string `json:"location,omitempty" validate:"required"`
	Title	 string `json:"title,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Balance float32 `json:"balance,omitempty"`
	OperatorCode string `json:"operatorcode,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type Wallet struct {
	Id primitive.ObjectID `json:"id,omitempty"`
	Method string `json:"method,omitempty"`
	OperatorCode string `json:"operatorcode,omitempty"`
	Providercode	 string `json:"providercode,omitempty"`
	Userid string `json:"userid,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Referenceid string `json:"referenceid,omitempty"`
	Action	 string `json:"action,omitempty"`
	Amount     float32 `json:"amount,omitempty"`	
	Signature string `json:"signature,omitempty"`
	Time string `json:"time,omitempty"`
	Status string `json:"status,omitempty"`
}