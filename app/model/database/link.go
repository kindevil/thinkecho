/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 23:16:17
 * @FilePath: /thinkecho/app/model/database/link.go
 * @Date: 2022-01-16 22:36:32
 * @Software: VS Code
 */
package database

import "gorm.io/gorm"

type Link struct {
	Lid         uint   `gorm:"column:lid;"`         // LID
	Name        string `gorm:"column:name;"`        // 链接名称
	Link        string `gorm:"column:link;"`        // 链接网址
	Description string `gorm:"column:description;"` // 链接描述
	Group       uint   `gorm:"column:group;"`       // 链接分组
}

/**
 * @description: 更新前操作
 * @param {*gorm.DB} tx
 * @return {*}
 */
func (l *Link) BeforeUpdate(tx *gorm.DB) (err error) {
	l.Lid = 0
	return
}

/**
 * @description: 添加链接
 * @param {*}
 * @return {*}
 */
func (l *Link) Create() uint {
	DB.Create(l)
	return l.Lid
}

/**
 * @description: 更新链接
 * @param {*}
 * @return {*}
 */
func (l *Link) Update() {
	DB.Where("lid = ?", l.Lid).Updates(l)
}

/**
 * @description: 删除链接
 * @param {int} lid
 * @return {*}
 */
func DeleteLink(lid int) {
	DB.Where("lid = ?", lid).Delete(&Link{})
}

/**
 * @description: 获取链接
 * @param {int} lid
 * @return {*}
 */
func GetLink(lid int) *Link {
	var link Link
	DB.Where("lid = ?", lid).First(&link)
	return &link
}

/**
 * @description: 获取所有链接
 * @param {*}
 * @return {*}
 */
func GetLinks() *[]Link {
	var links []Link
	DB.Find(&links)
	return &links
}
