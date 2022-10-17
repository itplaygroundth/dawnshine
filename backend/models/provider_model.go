package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Provider struct {
	Id primitive.ObjectID          `json:"id,omitempty"`
    ProviderCode string            `json:"providercode,omitempty" validate:"required"`
    OperatorCode string            `json:"operatorcode,omitempty" validate:"required"`
	SecretKey string 			   `json:"secretkey,omitempty"`
	AgentCurrency string           `json:"agentcurrency,omitempty" validate:"required"`
	BackendUrl string 			   `json:"backendurl,omitempty" validate:"required"`	
	UserName string                `json:"username,omitempty" validate:"required"`
	Password string                `json:"password,omitempty" validate:"required"`
}