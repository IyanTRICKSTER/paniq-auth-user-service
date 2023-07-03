package database

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"paniq-auth-user-service/pkg/entities"
	"testing"
)

func TestCreateDBConnection(t *testing.T) {

	//Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	var dbGlobal *Database

	t.Run("test create connection", func(t *testing.T) {
		db := Database{
			Host:     os.Getenv("DB_HOST"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_NAME"),
			DbPort:   os.Getenv("DB_PORT"),
		}

		err := db.Connect()
		assert.Nil(t, err)

		dbGlobal = &db

	})

	t.Run("test connection already established", func(t *testing.T) {
		err := dbGlobal.Connect()
		assert.Nil(t, err)
	})

	t.Run("test get active connection", func(t *testing.T) {
		assert.NotNil(t, dbGlobal.GetConnection())
	})

	t.Run("test can't establish connection using invalid host ip", func(t *testing.T) {
		db := Database{
			Host:     "0",
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_NAME"),
			DbPort:   os.Getenv("DB_PORT"),
		}

		err := db.Connect()
		assert.NotNil(t, err)
	})
}

func TestMigrationDB(t *testing.T) {

	//Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := Database{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		DbPort:   os.Getenv("DB_PORT"),
	}

	err := db.Connect()
	assert.Nil(t, err)

	t.Run("test migrate model", func(t *testing.T) {
		err := Migrate(db.GetConnection(), entities.UserEntity{})
		assert.Nil(t, err)
	})

	t.Run("test drop non defined table", func(t *testing.T) {
		err := DropTable(db.GetConnection(), entities.UserEntity{})
		assert.Nil(t, err)

		err = DropTable(db.GetConnection(), entities.UserEntity{})
		log.Println(err)
		assert.NotNil(t, err)

	})
}

func TestSeeder(t *testing.T) {

	//Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := Database{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		DbPort:   os.Getenv("DB_PORT"),
	}

	_ = db.Connect()

	_ = DropTable(db.GetConnection(), entities.UserEntity{})
	_ = DropTable(db.GetConnection(), entities.PermissionEntity{})
	_ = DropTable(db.GetConnection(), entities.RoleEntity{})
	_ = Migrate(db.GetConnection(), entities.UserEntity{})
	_ = Migrate(db.GetConnection(), entities.RoleEntity{})
	_ = Migrate(db.GetConnection(), entities.PermissionEntity{})

	t.Run("Test Run Permission Seeder", func(t *testing.T) {
		errs := PermissionSeeder(db.GetConnection())
		for _, err := range errs {
			assert.Nil(t, err)
		}
	})

	t.Run("Test Run Role Seeder", func(t *testing.T) {
		errs := RoleSeeder(db.GetConnection())
		for _, err := range errs {
			assert.Nil(t, err)
		}
	})

	t.Run("Test Run User Seeder", func(t *testing.T) {
		errs := UserSeeder(db.GetConnection())
		for _, err := range errs {
			assert.Nil(t, err)
		}
	})
}
