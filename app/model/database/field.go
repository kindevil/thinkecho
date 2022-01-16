/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 23:10:44
 * @FilePath: /thinkecho/app/model/database/field.go
 * @Date: 2022-01-16 22:35:48
 * @Software: VS Code
 */
package database

type Field struct {
	Cid        uint    `gorm:"column:cid;"`         // 内容ID
	Name       string  `gorm:"column:name;"`        // 名称
	Type       string  `gorm:"column:type;"`        // 类型
	StrValue   string  `gorm:"column:str_value;"`   // 字符串值
	IntValue   int     `gorm:"column:int_value;"`   // 整型值
	FloatValue float64 `gorm:"column:float_value;"` // 浮点型值
}

/**
 * @description: 添加字段
 * @param {*}
 * @return {int}
 */
func (f *Field) Create() {
	DB.Create(f)
}

/**
 * @description: 更新字段
 * @param {*}
 * @return {*}
 */
func (f *Field) Update() {
	DB.Where("cid = ? AND name = ?", f.Cid, f.Name).Updates(f)
}

/**
 * @description: 获取字段
 * @param {int} cid
 * @return {*}
 */
func GetFields(cid uint) *[]Field {
	var fields []Field
	DB.Where("cid = ?", cid).Find(&fields)
	return &fields
}

/**
 * @description: 删除字段
 * @param {int} cid
 * @param {string} name
 * @return {*}
 */
func DeleteField(cid uint, name string) {
	DB.Where("cid = ? AND name = ?", cid, name).Delete(&Field{})
}
