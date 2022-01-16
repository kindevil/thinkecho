/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:43:21
 * @FilePath: /thinkecho/app/service/archive/post.go
 * @Date: 2021-11-10 16:25:43
 * @Software: VS Code
 */
package archive

import (
	"strings"
	"thinkecho/app/database/content"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Archive
	Meta
}

func NewPost(c *gin.Context, content *content.Content) *Post {
	post := &Post{
		Archive{
			Cid:   content.Cid,
			Title: content.Title,
			Slug:  content.Slug,
			Author: Author{
				Name: content.Author.ScreenName,
				uid:  content.Author.UID,
			},
			allowComment: content.AllowComment,
			commentsNum:  content.CommentNum,
			text:         strings.ReplaceAll(content.Text, "<!--markdown-->", ""),
			password:     content.Password,
			created:      content.Created,
			fields:       content.Fields,
			context:      c,
			isPost:       true,
		},
		Meta{
			categories: content.Categories,
			tags:       content.Tags,
		},
	}

	return post
}
