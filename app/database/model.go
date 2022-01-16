/*
 * @Author: jia
 * @LastEditTime: 2021-11-16 09:56:08
 * @FilePath: /thinkecho/app/database/model.go
 * @Date: 2021-11-03 23:08:47
 * @Software: VS Code
 */
package database

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

/**
 * @description: 数据库初始化
 * @param {*}
 * @return {*}
 */
func InitDb() error {
	var err error
	var db *gorm.DB
	dbtype := viper.GetString("db_type")
	switch dbtype {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s", viper.GetString("mysql.user"), viper.GetString("mysql.password"), viper.GetString("mysql.host"), viper.GetString("mysql.port"), viper.GetString("mysql.db_name"), viper.GetString("mysql.charset"), viper.GetString("mysql.parse_time"), viper.GetString("mysql.local"))
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   viper.GetString("mysql.prefix"),
				SingularTable: false,
			},
		})
		if err != nil {
			return err
		}
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", viper.GetString("postgre.host"), viper.GetString("postgre.user"), viper.GetString("postgre.password"), viper.GetString("postgre.db_name"), viper.GetString("postgre.port"), viper.GetString("postgre.ssl_mode"), viper.GetString("postgre.timezone"))
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(viper.GetString("sqlite.file_path")), &gorm.Config{})
		if err != nil {
			return err
		}
	}

	// 是否为
	if viper.GetString("mode") == "debug" {
		DB = db.Debug()
	} else {
		DB = db
	}

	return nil
}
