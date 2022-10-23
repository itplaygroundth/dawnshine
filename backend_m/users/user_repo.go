package users
 

import (
	//"context"
	//"fmt"
	// "log"
	// "time"
	// "database/sql"
	//"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}
 
func ProvidUserRepository(DB *gorm.DB) UserRepository {
	return UserRepository{DB:DB}
}

func (u *UserRepository) FindAll() []User{
	var users []User
	u.DB.Find(&users)
	return users
}

func (u *UserRepository) FindByID(id uint64) User {
	var user User
	u.DB.First(&user,id)
	return user
}

func (u *UserRepository) Save(user User) User {
	u.DB.Save(&user)
	return user
}

func (u *UserRepository) Delete(user User){
	u.DB.Delete(&user)
}
// // type HandlerInterface interface {
// // 	configs *Handler
// // 	app *fiber.App

// // }
// func NewHandler(db *gorm.DB) *Handler {
// 	return &Handler{db:db}
// }
 

// func ConnectDB() *gorm.DB {

// 	//dsn := "root:12345678@tcp(localhost:3306)/go_basics?parseTime=true"
// 	dial := mysql.Open(DbConnection())
// 	db, err := gorm.Open(dial)
// 	if err != nil {
// 		panic(err)
// 	}

 
// 	fmt.Println("Connect to Mysql successfully")
// 	return db
// }