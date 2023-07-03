package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"paniq-auth-user-service/pkg/database"
	authControllers "paniq-auth-user-service/pkg/domains/auth/controllers"
	authUsecase "paniq-auth-user-service/pkg/domains/auth/usecases"
	"paniq-auth-user-service/pkg/domains/notification"
	userControllers "paniq-auth-user-service/pkg/domains/user/controllers"
	"paniq-auth-user-service/pkg/domains/user/repositories"
	userUsecase "paniq-auth-user-service/pkg/domains/user/usecases"
	"paniq-auth-user-service/pkg/entities"
	bcryptUtils "paniq-auth-user-service/pkg/utils/bcrypt"
	jwtUtils "paniq-auth-user-service/pkg/utils/jwt"
)

func main() {

	//Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Connect to the database
	db := database.Database{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		DbPort:   os.Getenv("DB_PORT"),
	}
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	//Migrate Entity
	_ = database.DropTable(db.GetConnection(), entities.UserEntity{})
	_ = database.DropTable(db.GetConnection(), entities.PermissionEntity{})
	_ = database.DropTable(db.GetConnection(), entities.RoleEntity{})
	_ = database.Migrate(db.GetConnection(), entities.UserEntity{})
	_ = database.Migrate(db.GetConnection(), entities.RoleEntity{})
	_ = database.Migrate(db.GetConnection(), entities.PermissionEntity{})

	//Run Seeders
	database.PermissionSeeder(db.GetConnection())
	database.RoleSeeder(db.GetConnection())
	database.UserSeeder(db.GetConnection())

	//Enable Gin Debugging Mode
	//gin.SetMode(gin.ReleaseMode)
	httpEngine := gin.Default()

	//JWT Service
	jwtSvc := jwtUtils.New()
	notificationSvc := notification.NewUsecase(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)

	//Setup Repository
	userRepo := repositories.NewUserRepo(db)
	//Setup Usecase
	authUsc := authUsecase.NewAuthUsecase(userRepo, jwtSvc, bcryptUtils.NewHashFunction())
	userUsc := userUsecase.NewUserUsecase(userRepo, bcryptUtils.NewHashFunction(), jwtSvc, notificationSvc)
	//Setup Controller
	authControllers.RunAuthController(httpEngine, authUsc)
	userControllers.RunUserController(httpEngine, userUsc)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "80"
	}

	err = httpEngine.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
