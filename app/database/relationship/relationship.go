/*
 * @Author: jia
 * @LastEditTime: 2021-11-16 10:01:47
 * @FilePath: /thinkecho/app/database/relationship/relationship.go
 * @Date: 2021-11-03 23:15:09
 * @Software: VS Code
 */
package relationship

import (
	"thinkecho/app/database"

	"github.com/spf13/viper"
)

type Relationship struct {
	Cid uint `gorm:"column:cid;"` // 内容ID
	Mid uint `gorm:"column:mid;"` // Meta ID
}

/**
 * @description: 定义表名
 * @param {*}
 * @return {*}
 */
func (Relationship) TableName() string {
	return viper.GetString("mysql.prefix") + "relationships"
}

/**
 * @description: 添加关系
 * @param {*}
 * @return {*}
 */
func (r *Relationship) Create() {
	database.DB.Create(r)
}

/**
 * @description: 删除关系
 * @param {int} cid
 * @return {*}
 */
func DeleteRelationship(cid int) {
	database.DB.Where("cid = ?", cid).Delete(&Relationship{})
}
