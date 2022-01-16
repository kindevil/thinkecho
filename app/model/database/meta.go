/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 23:15:20
 * @FilePath: /thinkecho/app/model/database/meta.go
 * @Date: 2022-01-16 22:37:10
 * @Software: VS Code
 */
package database

import (
	"reflect"

	"gorm.io/gorm"
)

type Meta struct {
	Mid         uint   `gorm:"column:mid;"`         // Mid
	Name        string `gorm:"column:name;"`        // 名称
	Slug        string `gorm:"column:slug;"`        // 缩略名
	Type        string `gorm:"column:type;"`        // 类型
	Description string `gorm:"column:description;"` // 描述
	Count       uint   `gorm:"column:count;"`       // 计数
	Order       uint   `gorm:"column:order;"`       // 排序
	Parent      uint   `gorm:"column:parent;"`      // 父
}

/**
 * @description: 更新前操作
 * @param {*gorm.DB} tx
 * @return {*}
 */
func (m *Meta) BeforeUpdate(tx *gorm.DB) (err error) {
	m.Mid = 0
	return
}

/**
 * @description: 添加Meta
 * @param {*}
 * @return {*}
 */
func (m *Meta) Create() uint {
	DB.Create(m)
	return m.Mid
}

/**
 * @description: 更新Meta
 * @param {*}
 * @return {*}
 */
func (m *Meta) Update() {
	DB.Where("mid = ?", m.Mid).Updates(m)
}

/**
 * @description: 判断结构体是否为空
 * @param {*}
 * @return {*}
 */
func (c *Meta) IsEmpty() bool {
	return reflect.DeepEqual(c, Meta{})
}

/**
 * @description: 删除Meta
 * @param {int} mid
 * @return {*}
 */
func DeleteMeta(mid int) {
	DB.Where("mid = ?", mid).Delete(&Meta{})
}

/**
 * @description: 获取metas
 * @param {string} slug
 * @param {string} metaType
 * @return {*}
 */
func GetMeta(slug string, metaType string) *Meta {
	var meta Meta
	DB.Where("slug = ? AND type = ?", slug, metaType).Preload("Relationship").First(&meta)
	return &meta
}
