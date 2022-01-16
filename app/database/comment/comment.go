/*
 * @Author: jia
 * @LastEditTime: 2021-11-16 10:01:52
 * @FilePath: /thinkecho/app/database/comment/comment.go
 * @Date: 2021-11-03 23:14:30
 * @Software: VS Code
 */
package comment

import "github.com/spf13/viper"

type Comment struct {
	Coid     int    `gorm:"column:coid;"`      // 评论ID
	Cid      int    `gorm:"column:cid;"`       // 内容ID
	Created  int    `gorm:"column:created;"`   // 创建时间
	Author   string `gorm:"column:author;"`    // 作者
	AuthorId int    `gorm:"column:author_id;"` // 作者ID
	OwnerId  int    `gorm:"column:owner_id;"`  // 所有者ID
	Mail     string `gorm:"column:mail;"`      // Email
	URL      string `gorm:"column:url;"`       // 网址
	IP       string `gorm:"column:ip;"`        // IP地址
	Agent    string `gorm:"column:agent;"`     // Agent信息
	Text     string `gorm:"column:text;"`      // 评论内容
	Type     string `gorm:"column:type;"`      // 评论类型
	Status   string `gorm:"column:status;"`    // 评论状态
	Parent   int    `gorm:"column:parent;"`    // 父评论
}

/**
 * @description: 定义表名
 * @param {*}
 * @return {*}
 */
func (Comment) TableName() string {
	return viper.GetString("mysql.prefix") + "comments"
}
