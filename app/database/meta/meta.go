/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 23:16:53
 * @FilePath: /thinkecho/app/database/meta/meta.go
 * @Date: 2021-11-03 23:14:47
 * @Software: VS Code
 */
package meta

import (
	"reflect"
	"thinkecho/app/database"
	"thinkecho/app/database/relationship"

	"github.com/spf13/viper"
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

	Relationship []relationship.Relationship `gorm:"foreignKey:mid;joinForeignKey:mid;references:mid;joinReferences:mid;"`
}

/**
 * @description: 定义表名
 * @param {*}
 * @return {*}
 */
func (Meta) TableName() string {
	return viper.GetString("mysql.prefix") + "metas"
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
	database.DB.Create(m)
	return m.Mid
}

/**
 * @description: 更新Meta
 * @param {*}
 * @return {*}
 */
func (m *Meta) Update() {
	database.DB.Where("mid = ?", m.Mid).Updates(m)
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
	database.DB.Where("mid = ?", mid).Delete(&Meta{})
}

/**
 * @description: 获取metas
 * @param {string} slug
 * @param {string} metaType
 * @return {*}
 */
func GetMeta(slug string, metaType string) *Meta {
	var meta Meta
	database.DB.Where("slug = ? AND type = ?", slug, metaType).Preload("Relationship").First(&meta)
	return &meta
}
