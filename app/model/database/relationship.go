/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 23:17:53
 * @FilePath: /thinkecho/app/model/database/relationship.go
 * @Date: 2022-01-16 22:38:25
 * @Software: VS Code
 */
package database

type Relationship struct {
	Cid uint `gorm:"column:cid;"` // 内容ID
	Mid uint `gorm:"column:mid;"` // Meta ID
}

/**
 * @description: 添加关系
 * @param {*}
 * @return {*}
 */
func (r *Relationship) Create() {
	DB.Create(r)
}

/**
 * @description: 删除关系
 * @param {int} cid
 * @return {*}
 */
func DeleteRelationship(cid int) {
	DB.Where("cid = ?", cid).Delete(&Relationship{})
}
