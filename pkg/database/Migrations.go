package database

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"paniq-auth-user-service/pkg/entities"
)

func Migrate[Entity entities.UserEntity | entities.PermissionEntity | entities.RoleEntity](conn *gorm.DB, entity Entity) error {
	err := conn.AutoMigrate(&entity)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func DropTable[Entity entities.UserEntity | entities.PermissionEntity | entities.RoleEntity](conn *gorm.DB, entity Entity) error {

	if !conn.Migrator().HasTable(entity) {
		return errors.New("table is not exists")
	}

	err := conn.Migrator().DropTable(entity)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
