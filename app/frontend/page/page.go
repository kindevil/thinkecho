/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 22:49:18
 * @FilePath: /thinkecho/app/frontend/page/page.go
 * @Date: 2021-11-06 11:32:31
 * @Software: VS Code
 */
package page

import (
	"net/http"
	"thinkecho/app/database/content"
	"thinkecho/app/pkg/theme"
	"thinkecho/app/service"
	"thinkecho/app/service/archive"

	"github.com/gin-gonic/gin"
)

func Page(c *gin.Context) {
	slug := c.Param("slug")

	data := content.GetContent(slug, "page", "")
	if data.IsEmpty() || data.Cid == 0 {
		service.Error404(c)
		return
	}

	c.Set("archiveTitle", data.Title)
	c.Set("pageType", "page")
	c.Set("pageSlug", data.Slug)
	c.Set("cid", data.Cid)

	page := archive.NewPage(c, data)

	theme.HTML(c, http.StatusOK, "page.tmpl", gin.H{"archive": page})
}
