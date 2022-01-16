/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 23:16:47
 * @FilePath: /thinkecho/app/model/database/content.go
 * @Date: 2022-01-16 22:34:20
 * @Software: VS Code
 */
package database

import "gorm.io/gorm"

//Content 内容
type Content struct {
	Cid          uint   `gorm:"column:cid;"`          // 内容ID
	Title        string `gorm:"column:title;"`        // 标题
	Slug         string `gorm:"column:slug;"`         // 缩略名
	Created      int64  `gorm:"column:created;"`      // 创建日期
	Modified     int64  `gorm:"column:modified;"`     // 修改日期
	Text         string `gorm:"column:text;"`         // 内容
	Order        uint   `gorm:"column:order;"`        // 排序
	AuthorId     uint   `gorm:"column:authorId;"`     // 作者ID
	Template     uint   `gorm:"column:template;"`     // 模版
	Type         string `gorm:"column:type;"`         // 类型
	Status       string `gorm:"column:status;"`       // 状态
	Password     string `gorm:"column:password;"`     // 密码
	CommentNum   uint   `gorm:"column:commentsNum;"`  // 评论数
	AllowComment uint   `gorm:"column:allowComment;"` // 是否允许评论
	AllowPing    uint   `gorm:"column:allowPing;"`    // 是否允许ping
	AllowFeed    uint   `gorm:"column:allowFeed;"`    // 是否允许feed
	Parent       uint   `gorm:"column:parent;"`       // 父用户ID
}

/**
 * @description: 更新前操作
 * @param {*gorm.DB} tx
 * @return {*}
 */
func (c *Content) BeforeUpdate(tx *gorm.DB) (err error) {
	c.Cid = 0
	return
}

/**
 * @description: 添加内容
 * @param {*}
 * @return {*}
 */
func (c *Content) Create() uint {
	DB.Create(c)
	return c.Cid
}

/**
 * @description: 更新内容
 * @param {*}
 * @return {*}
 */
func (c *Content) Update() {
	DB.Where("cid = ?", c.Cid).Updates(c)
}

/**
 * @description:
 * @param {string} contentType
 * @param {[]int} cids
 * @param {string} status
 * @param {uint} authorId
 * @param {string} keywords
 * @param {int} limit
 * @param {int} offset
 * @param {string} order
 * @return {*}
 */
func GetContents(contentType string, cids []int, status string, authorId uint, keywords string, limit int, offset int, order string) *[]Content {
	var contents []Content
	db := DB

	if contentType != "" {
		db = db.Where("type = ?", contentType)
	}

	if len(cids) > 0 {
		db = db.Where("cid IN ?", cids)
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if authorId != 0 {
		db = db.Where("authorId = ?", authorId)
	}

	if keywords != "" {
		db = db.Where("title LIKE ? OR text LIKE ? ", "%"+keywords+"%", "%"+keywords+"%")
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	switch order {
	case "desc":
		db = db.Order("Cid desc")
	case "asc":
		db = db.Order("Cid asc")
	}

	db.Find(&contents)

	return &contents
}

/**
 * @description:
 * @param {string} contentType
 * @param {[]int} cids
 * @param {string} status
 * @param {uint} authorId
 * @param {string} keywords
 * @return {*}
 */
func GetContentCount(contentType string, cids []int, status string, authorId uint, keywords string) int64 {
	var contents []*Content

	db := DB.Select("cid")
	if contentType != "" {
		db = db.Where("type = ?", contentType)
	}

	if len(cids) > 0 {
		db = db.Where("cid IN ?", cids)
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if authorId != 0 {
		db = db.Where("authorId = ?", authorId)
	}

	if keywords != "" {
		db = db.Where("title LIKE ? OR text LIKE ? ", "%"+keywords+"%", "%"+keywords+"%")
	}

	return db.Find(&contents).RowsAffected
}

/**
 * @description: 添加内容评论计数
 * @param {int} cid
 * @return {*}
 */
func AddCommentCount(cid int) {
	var content Content
	DB.Select("commentsNum").Where("cid = ?", cid).First(&content)
	content.CommentNum++
	DB.Model(&Content{}).Where("cid = ?", cid).Update("commentNum", content.CommentNum)
}

/**
 * @description:减少评论计数
 * @param {int} cid
 * @return {*}
 */
func MinusCommentCount(cid int) {
	var content Content
	DB.Select("commentsNum").Where("cid = ?", cid).First(&content)
	content.CommentNum--
	DB.Model(&Content{}).Where("cid = ?", cid).Update("commentNum", content.CommentNum)
}

/**
 * @description: 判断内容是否存在
 * @param {int} cid
 * @param {int} slug
 * @return {*}
 */
func ContentExist(cid int, slug string) bool {
	var content Content
	if cid != 0 {
		return DB.Where("cid == ?", cid).First(&content).RowsAffected > 0
	}
	return DB.Where("slug == ?", slug).First(&content).RowsAffected > 0
}

/**
 * @description: 通过cid或者slug获取内容
 * @param {*}
 * @return {*}
 */
func GetContent(slug string, contentType string, status string) *Content {
	var content Content

	db := DB.Preload("Author").Preload("Fields")
	if contentType != "" {
		db = db.Where("type = ?", contentType)
	}

	if contentType == "post" {
		db = db.Preload("Categories", "type = ?", "category").Preload("Tags", "type = ?", "tag")
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	db.Where("slug = ?", slug).Find(&content)

	return &content
}

/**
 * @description: 删除内容
 * @param {int} cid
 * @param {string} slug
 * @return {*}
 */
func DeleteContent(cid int, slug string) {
	if cid != 0 {
		DB.Where("cid = ?", cid).Delete(&Content{})
	} else {
		DB.Where("slug = ?", slug).Delete(&Content{})
	}
}

/**
 * @description: 获取页面标题
 * @param {*}
 * @return {*}
 */
func GetPageTitle() *[]Content {
	var contents []Content
	DB.Select("title,slug").Where("type = ? AND status = ?", "page", "publish").Order("`order` asc").Find(&contents)
	return &contents
}

/**
 * @description: 上一篇文章
 * @param {int64} created
 * @return {*}
 */
func GetContentPrev(created int64) *Content {
	var content Content
	DB.Where("created < ? AND status = ? AND type = ? AND password IS NULL OR password = ''", created, "publish", "post").Order("`created` desc").First(&content)
	return &content
}

/**
 * @description: 下一篇文章
 * @param {int64} created
 * @return {*}
 */
func GetContentNext(created int64) *Content {
	var content Content
	DB.Where("created > ? AND status = ? AND type = ? AND password IS NULL OR password = ''", created, "publish", "post").Order("`created` desc").First(&content)
	return &content
}
