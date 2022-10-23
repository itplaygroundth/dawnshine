package main

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"cryptjoshi/users"
)

func initUserAPI(db *gorm.DB) users.UserAPI {
	wire.Build(users.ProvidUserRepository,users.ProvideUserService,users.ProvideUserAPI)
	return users.UserAPI{}
}

