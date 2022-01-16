/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 23:18:23
 * @FilePath: /thinkecho/app/model/database/user.go
 * @Date: 2022-01-16 22:39:14
 * @Software: VS Code
 */
package database

import (
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
	DB.Create(u)
	return u.UID
}

/**
 * @description: 更新用户信息
 * @param {*}
 * @return {*}
 */
func (u *User) Update() {
	DB.Model(&User{}).Updates(u)
}

/**
 * @description: 更新用户密码
 * @param {*}
 * @return {*}
 */
func UpdateUserPassword(uid int, password string) {
	DB.Model(&User{}).Where("uid = ?", uid).Update("password", password)
}

/**
 * @description: 通过uid获取用户信息
 * @param {int} uid
 * @return {*}
 */
func GetUserByID(uid int) *User {
	var user *User
	DB.Where("uid = ?", uid).First(&user)
	return user
}

/**
 * @description: 通过name获取用户信息
 * @param {string} name
 * @return {*}
 */
func GetUserByName(name string) *User {
	var user *User
	DB.Where("name = ?", name).First(&user)
	return user
}
