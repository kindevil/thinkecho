/*
 * @Author: jia
 * @LastEditTime: 2021-11-16 10:01:44
 * @FilePath: /thinkecho/app/database/user/user.go
 * @Date: 2021-11-03 23:15:17
 * @Software: VS Code
 */
package user

import (
	"reflect"
	"thinkecho/app/database"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type User struct {
	UID        uint   `gorm:"column:uid;"`        // UID
	Name       string `gorm:"column:name;"`       // 用户名
	Password   string `gorm:"column:password;"`   // 密码
	Mail       string `gorm:"column:mail;"`       // 邮箱
	URL        string `gorm:"column:url;"`        // 网址
	ScreenName string `gorm:"column:screenName;"` // 昵称
	Created    uint   `gorm:"column:created;"`    // 创建时间
	Activated  uint   `gorm:"column:activated;"`  // 是否激活
	Logged     uint   `gorm:"column:logged;"`     // 登录时间
	Group      string `gorm:"column:group;"`      // 用户组
	AuthCode   string `gorm:"column:authcode;"`   // authcode
}

/**
 * @description: 定义表名
 * @param {*}
 * @return {*}
 */
func (User) TableName() string {
	return viper.GetString("mysql.prefix") + "users"
}

/**
 * @description: 更新前操作
 * @param {*gorm.DB} tx
 * @return {*}
 */
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UID = 0
	return
}

/**
 * @description: 添加用户
 * @param {*}
 * @return {int}
 */
func (u *User) Create() uint {
	database.DB.Create(u)
	return u.UID
}

/**
 * @description: 更新用户信息
 * @param {*}
 * @return {*}
 */
func (u *User) Update() {
	database.DB.Model(&User{}).Updates(u)
}

/**
 * @description: 判断结构体是否为空
 * @param {*}
 * @return {*}
 */
func (c *User) IsEmpty() bool {
	return reflect.DeepEqual(c, User{})
}

/**
 * @description: 更新用户密码
 * @param {*}
 * @return {*}
 */
func UpdateUserPassword(uid int, password string) {
	database.DB.Model(&User{}).Where("uid = ?", uid).Update("password", password)
}

/**
 * @description: 通过uid获取用户信息
 * @param {int} uid
 * @return {*}
 */
func GetUserByID(uid int) *User {
	var user *User
	database.DB.Where("uid = ?", uid).First(&user)
	return user
}

/**
 * @description: 通过name获取用户信息
 * @param {string} name
 * @return {*}
 */
func GetUserByName(name string) *User {
	var user *User
	database.DB.Where("name = ?", name).First(&user)
	return user
}
