package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	Host     string
	Username string
	Password string
	DbName   string
	DbPort   string
	conn     *gorm.DB
}

func (d *Database) Connect() error {
	if d.conn == nil {
		dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", d.Username, d.Password, d.Host, d.DbPort, d.DbName)
		conn, err := gorm.Open(mysql.Open(dbURI), &gorm.Config{})
		if err != nil {
			fmt.Println("Can't establish a connection: ", err)
			return err
		}
		fmt.Println("Connection has been established")
		d.conn = conn
	} else {
		fmt.Println("Connection already established")
	}

	return nil
}

func (d *Database) GetConnection() *gorm.DB {
	return d.conn
}
