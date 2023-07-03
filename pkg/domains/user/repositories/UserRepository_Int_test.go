package repositories

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/contracts/statusCodes"
	"paniq-auth-user-service/pkg/database"
	"paniq-auth-user-service/pkg/entities"
	bcryptUtils "paniq-auth-user-service/pkg/utils/bcrypt"
	"testing"
	"time"
)

var userRepo contracts.IUserRepository
var hashFunction contracts.IHash

func init() {
	//Load .env file
	if err := godotenv.Load("../../../../.env"); err != nil {
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

	userRepo = NewUserRepo(db)

	hashFunction = bcryptUtils.NewHashFunction()
}

func TestUserRepository_FetchAllUser(t *testing.T) {
	t.Run("Fetch Success", func(t *testing.T) {

		ctx := context.Background()

		res := userRepo.FetchAllUser(ctx)
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.NotNil(t, res.GetData())
	})
}

func TestUserRepository_FetchUserByEmail(t *testing.T) {

	t.Run("Fetch Success", func(t *testing.T) {

		res := userRepo.FetchUserByEmail(context.Background(), "iyanpratama2002@gmail.com")
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.NotNil(t, res.GetData())
	})

	t.Run("User Not Found", func(t *testing.T) {
		res := userRepo.FetchUserByEmail(context.Background(), "iyanpratama@gmail.com")
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ModelNotFound, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("Fetch Error", func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		res := userRepo.FetchUserByEmail(ctx, "iyanpratama@gmail.com")
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})
}

func TestUserRepository_FetchUserByID(t *testing.T) {

	t.Run("Fetch Success", func(t *testing.T) {
		res := userRepo.FetchUserByID(context.Background(), 1)
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.NotNil(t, res.GetData())
	})

	t.Run("User Not Found", func(t *testing.T) {
		res := userRepo.FetchUserByID(context.Background(), 1000)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ModelNotFound, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("Fetch Error", func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		res := userRepo.FetchUserByID(ctx, 1)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

}

func TestUserRepository_CreateUser(t *testing.T) {

	netRandom := time.Now().UTC().UnixNano()

	t.Run("Create Success", func(t *testing.T) {

		user := entities.UserEntity{
			Role:       entities.RoleEntity{},
			RoleID:     uint(3),
			Username:   fmt.Sprintf("%v-%v", "iyan", netRandom),
			Password:   hashFunction.Hash("iyan12345"),
			Email:      fmt.Sprintf("%v-%v@gmail.com", "iyan", netRandom),
			Avatar:     "https://i.pravatar.cc/300",
			NIM:        nil,
			NIP:        nil,
			Major:      "CS Degree",
			ResetToken: "",
		}

		res := userRepo.CreateUser(context.Background(), user)

		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.NotNil(t, res.GetData().(entities.UserEntity).CreatedAt)
	})

	t.Run("Create but duplicated", func(t *testing.T) {

		user := entities.UserEntity{
			Role:       entities.RoleEntity{},
			RoleID:     uint(3),
			Username:   "akiyan2002",
			Password:   hashFunction.Hash("iyan12345"),
			Email:      "akiyan2002@gmail.com",
			Avatar:     "https://i.pravatar.cc/300",
			NIM:        nil,
			NIP:        nil,
			Major:      "CS Degree",
			ResetToken: "",
		}

		res := userRepo.CreateUser(context.Background(), user)

		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

}

func TestUserRepository_CreateMultiUsers(t *testing.T) {

	t.Run("Create Success", func(t *testing.T) {

		var users []entities.UserEntity

		for i := 0; i < 200; i++ {
			user := entities.UserEntity{
				Role:       entities.RoleEntity{},
				RoleID:     uint(3),
				Username:   fmt.Sprintf("%v-%v", "iyan", i),
				Password:   hashFunction.Hash("iyan12345"),
				Email:      fmt.Sprintf("%v-%v@gmail.com", "iyan", i),
				Avatar:     "https://i.pravatar.cc/300",
				NIM:        nil,
				NIP:        nil,
				Major:      "CS Degree",
				ResetToken: "",
			}
			users = append(users, user)
		}

		res := userRepo.CreateUsers(context.Background(), users)

		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("Create Fail contain duplicate", func(t *testing.T) {

		var users []entities.UserEntity

		for i := 0; i < 1; i++ {
			user := entities.UserEntity{
				Role:       entities.RoleEntity{},
				RoleID:     uint(3),
				Username:   fmt.Sprintf("%v-%v", "iyan", i),
				Password:   hashFunction.Hash("iyan12345"),
				Email:      fmt.Sprintf("%v-%v@gmail.com", "iyan", i),
				Avatar:     "https://i.pravatar.cc/300",
				NIM:        nil,
				NIP:        nil,
				Major:      "CS Degree",
				ResetToken: "",
			}
			users = append(users, user)
		}

		res := userRepo.CreateUsers(context.Background(), users)

		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

}

func TestUserRepository_UpdateUser(t *testing.T) {

	t.Run("Update Success", func(t *testing.T) {

		user := entities.UserEntity{
			ID:         1,
			Role:       entities.RoleEntity{},
			RoleID:     uint(3),
			Username:   fmt.Sprintf("%v-%v", "iyan", 669966),
			Password:   hashFunction.Hash("iyan12345"),
			Email:      fmt.Sprintf("%v-%v@gmail.com", "iyan", 669966),
			Avatar:     "https://i.pravatar.cc/300",
			NIM:        nil,
			NIP:        nil,
			Major:      "CS Degree",
			ResetToken: "",
		}

		res := userRepo.UpdateUser(context.Background(), user)

		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.NotNil(t, res.GetData().(entities.UserEntity).CreatedAt)
	})

	t.Run("Update Failed context canceled", func(t *testing.T) {

		user := entities.UserEntity{
			ID:         1,
			Role:       entities.RoleEntity{},
			RoleID:     uint(3),
			Username:   fmt.Sprintf("%v-%v", "iyan", 669966),
			Password:   hashFunction.Hash("iyan12345"),
			Email:      fmt.Sprintf("%v-%v@gmail.com", "iyan", 669966),
			Avatar:     "https://i.pravatar.cc/300",
			NIM:        nil,
			NIP:        nil,
			Major:      "CS Degree",
			ResetToken: "",
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		res := userRepo.UpdateUser(ctx, user)

		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

}

func TestUserRepository_FetchUserACL(t *testing.T) {

	t.Run("fetch user ACL success", func(t *testing.T) {
		res := userRepo.FetchUserACL(context.Background(), 1, apiResources.USER)
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.NotEmpty(t, res.GetData())
		log.Println(res.GetData())
	})

	t.Run("fetch user ACL failed resource not exists", func(t *testing.T) {
		res := userRepo.FetchUserACL(context.Background(), 1, "WKWKWK")
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.Empty(t, res.GetData().(entities.UserEntity).Role.Permissions)
	})

	t.Run("fetch user ACL failed context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		res := userRepo.FetchUserACL(ctx, 1, "WKWKWK")
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

}
