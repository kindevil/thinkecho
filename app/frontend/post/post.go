/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 22:58:01
 * @FilePath: /thinkecho/app/frontend/post/post.go
 * @Date: 2021-11-10 16:36:04
 * @Software: VS Code
 */
package post

import (
	"fmt"
	"net/http"
	"thinkecho/app/database/content"
	"thinkecho/app/pkg/theme"
	"thinkecho/app/service"
	"thinkecho/app/service/archive"

	"github.com/gin-gonic/gin"
)

func Post(c *gin.Context) {
	slug := c.Param("slug")

	fmt.Println(slug)

	data := content.GetContent(slug, "post", "publish")
	if data.IsEmpty() || data.Cid == 0 {
		service.Error404(c)
		return
	}

	c.Set("archiveTitle", data.Title)
	c.Set("pageType", "post")
	c.Set("pageSlug", data.Slug)

	post := archive.NewPost(c, data)

	theme.HTML(c, http.StatusOK, "post.tmpl", gin.H{
		"archive": post,
	})
}
