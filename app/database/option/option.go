/*
 * @Author: jia
 * @LastEditTime: 2021-11-16 10:01:40
 * @FilePath: /thinkecho/app/database/option/option.go
 * @Date: 2021-11-03 23:14:55
 * @Software: VS Code
 */
package option

import (
	"thinkecho/app/database"

	"github.com/spf13/viper"
)

type Option struct {
	Name  string `gorm:"column:name;"`  // 配置名
	User  string `gorm:"column:user;"`  // 配置用户
	Value string `gorm:"column:value;"` // 配置值
}

/**
 * @description: 定义表名
 * @param {*}
 * @return {*}
 */
func (Option) TableName() string {
	return viper.GetString("mysql.prefix") + "options"
}

/**
 * @description: 添加配置
 * @param {*}
 * @return {*}
 */
func (o *Option) Create() {
	database.DB.Create(o)
}

/**
 * @description:
 * @param {*}
 * @return {*}
 */
func (o *Option) Update() {
	database.DB.Where("name = ?", o.Name).Updates(o)
}

/**
 * @description: 删除配置
 * @param {string} name
 * @return {*}
 */
func DeleteOption(name string) {
	database.DB.Where("name = ?", name).Delete(&Option{})
}

/**
 * @description: 获取配置
 * @param {string} name
 * @return {*}
 */
func GetOption(name string) *Option {
	var option Option
	database.DB.Where("name = ?", name).First(&option)
	return &option
}

/**
 * @description: 获取所有配置
 * @param {*}
 * @return {*}
 */
func GetOptions() *[]Option {
	var options []Option
	database.DB.Find(&options)
	return &options
}
