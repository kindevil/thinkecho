/*
 * @Author: jia
 * @LastEditTime: 2021-11-16 10:01:24
 * @FilePath: /thinkecho/app/database/field/field.go
 * @Date: 2021-11-03 23:14:12
 * @Software: VS Code
 */
package field

import (
	"thinkecho/app/database"

	"github.com/spf13/viper"
)

type Field struct {
	Cid        uint    `gorm:"column:cid;"`         // 内容ID
	Name       string  `gorm:"column:name;"`        // 名称
	Type       string  `gorm:"column:type;"`        // 类型
	StrValue   string  `gorm:"column:str_value;"`   // 字符串值
	IntValue   int     `gorm:"column:int_value;"`   // 整型值
	FloatValue float64 `gorm:"column:float_value;"` // 浮点型值
}

/**
 * @description: 定义表名
 * @param {*}
 * @return {*}
 */
func (Field) TableName() string {
	return viper.GetString("mysql.prefix") + "fields"
}

/**
 * @description: 添加字段
 * @param {*}
 * @return {int}
 */
func (f *Field) Create() {
	database.DB.Create(f)
}

/**
 * @description: 更新字段
 * @param {*}
 * @return {*}
 */
func (f *Field) Update() {
	database.DB.Where("cid = ? AND name = ?", f.Cid, f.Name).Updates(f)
}

/**
 * @description: 获取字段
 * @param {int} cid
 * @return {*}
 */
func GetFields(cid uint) *[]Field {
	var fields []Field
	database.DB.Where("cid = ?", cid).Find(&fields)
	return &fields
}

/**
 * @description: 删除字段
 * @param {int} cid
 * @param {string} name
 * @return {*}
 */
func DeleteField(cid uint, name string) {
	database.DB.Where("cid = ? AND name = ?", cid, name).Delete(&Field{})
}
