/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 22:35:03
 * @FilePath: /thinkecho/app/model/database/comment.go
 * @Date: 2022-01-16 22:35:02
 * @Software: VS Code
 */
package database

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
