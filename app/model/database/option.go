/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 23:17:24
 * @FilePath: /thinkecho/app/model/database/option.go
 * @Date: 2022-01-16 22:37:51
 * @Software: VS Code
 */
package database

type Option struct {
	Name  string `gorm:"column:name;"`  // 配置名
	User  string `gorm:"column:user;"`  // 配置用户
	Value string `gorm:"column:value;"` // 配置值
}

/**
 * @description: 添加配置
 * @param {*}
 * @return {*}
 */
func (o *Option) Create() {
	DB.Create(o)
}

/**
 * @description:
 * @param {*}
 * @return {*}
 */
func (o *Option) Update() {
	DB.Where("name = ?", o.Name).Updates(o)
}

/**
 * @description: 删除配置
 * @param {string} name
 * @return {*}
 */
func DeleteOption(name string) {
	DB.Where("name = ?", name).Delete(&Option{})
}

/**
 * @description: 获取配置
 * @param {string} name
 * @return {*}
 */
func GetOption(name string) *Option {
	var option Option
	DB.Where("name = ?", name).First(&option)
	return &option
}

/**
 * @description: 获取所有配置
 * @param {*}
 * @return {*}
 */
func GetOptions() *[]Option {
	var options []Option
	DB.Find(&options)
	return &options
}
