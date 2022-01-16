/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 22:23:03
 * @FilePath: /thinkecho/app/pkg/config/config.go
 * @Date: 2022-01-16 22:23:03
 * @Software: VS Code
 */
package config

import (
	"bytes"
	"fmt"
	"os"
	"thinkecho/app/pkg/file"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

type Config struct {
	Address string `mapstructure:"address" toml:"address"` // http服务器监听端口
	Mode    string `mapstructure:"release" toml:"release"` // 当前模式,debug模式和release模式
	DBType  string `mapstructure:"db_type" toml:"db_type"` // 数据库类型sqlite,mysql,postgresql
	MySQL   *Mysql
	SQLite  *SQLite
	Postgre *PostgreSQL
}

//Mysql 数据库
type Mysql struct {
	Host      string `mapstructure:"host" toml:"host"`            // 主机名
	Port      string `mapstructure:"port" toml:"port"`            // 端口
	User      string `mapstructure:"user" toml:"user"`            // 用户名
	Password  string `mapstructure:"password" toml:"password"`    // 密码
	DBName    string `mapstructure:"dbname" toml:"db_name"`       // 数据库密码
	Prefix    string `mapstructure:"prefix" toml:"prefix"`        // 表前缀
	Charset   string `mapstructure:"charset" toml:"charset"`      // 编码格式
	ParseTime string `mapstructure:"parseTime" toml:"parse_time"` // 把数据库datetime和date类型转换为golang的time.Time类型
	Local     string `mapstructure:"loc" toml:"local"`            // 时区
}

//SQLite 数据库
type SQLite struct {
	FilePath string `mapstructure:"filepath" toml:"file_path"` // db文件位置路径
}

//PostgreSQL 数据库
type PostgreSQL struct {
	Host     string `mapstructure:"host" toml:"host"`          // 主机名
	Port     string `mapstructure:"port" toml:"port"`          // 端口
	User     string `mapstructure:"user" toml:"user"`          // 用户名
	Password string `mapstructure:"password" toml:"password"`  // 密码
	DBName   string `mapstructure:"dbname" toml:"db_name"`     // 数据库密码
	SSLmode  string `mapstructure:"ssl_mode" toml:"ssl_mode"`  // 是否打开ssl模式
	TimeZone string `mapstructure:"time_zone" toml:"timezone"` // 时区
}

var configFile = "config.toml"

/**
 * @description: 初始化配置文件
 * @param {*}
 * @return {*}
 */
func InitConfig() error {
	var err error
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	if !file.PathExists(configFile) {
		data := initData()

		if err = viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
			return fmt.Errorf("读取默认配置文件 %s 内容失败: %s", configFile, err)
		}

		if err = viper.SafeWriteConfigAs(configFile); err != nil {
			if os.IsNotExist(err) {
				err = viper.WriteConfigAs(configFile)
				if err != nil {
					return fmt.Errorf("保存默认配置到文件 %s 失败: %s", configFile, err)
				}
			}
		}
	}

	if err = viper.ReadInConfig(); err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

/**
 * @description: 配置文件初始化数据
 * @param {*}
 * @return {*}
 */
func initData() []byte {
	mysql := &Mysql{
		Host:      "localhost",
		Port:      "3306",
		User:      "thinkecho",
		Password:  "passwd",
		DBName:    "thinkecho",
		Charset:   "utf8mb4",
		ParseTime: "True",
		Local:     "Asia/Shanghai",
	}

	postgreSQL := &PostgreSQL{
		Host:     "localhost",
		Port:     "9920",
		User:     "thinkecho",
		Password: "passwd",
		DBName:   "thinkecho",
		SSLmode:  "disable",
		TimeZone: "Asia/Shanghai",
	}

	sqlite := &SQLite{
		FilePath: "thinkecho.db",
	}

	config := &Config{
		Address: ":3000",
		Mode:    "debug",
		DBType:  "mysql",
		MySQL:   mysql,
		Postgre: postgreSQL,
		SQLite:  sqlite,
	}

	tomlByte, err := toml.Marshal(config)
	if err != nil {
		panic(err)
	}

	return tomlByte
}
