/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:43:45
 * @FilePath: /thinkecho/app/service/archive/page.go
 * @Date: 2021-11-10 16:25:34
 * @Software: VS Code
 */
package archive

import (
	"strings"
	"thinkecho/app/database/content"

	"github.com/gin-gonic/gin"
)

func NewPage(c *gin.Context, content *content.Content) *Archive {
	archive := &Archive{
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
		isPost:       false,
	}

	return archive
}
