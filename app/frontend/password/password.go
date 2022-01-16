/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:55:05
 * @FilePath: /thinkecho/app/frontend/password/password.go
 * @Date: 2021-11-15 14:50:20
 * @Software: VS Code
 */
package password

import (
	"fmt"
	"net/http"
	"thinkecho/app/database/content"
	"thinkecho/app/service"

	"github.com/gin-gonic/gin"
)

func Password(c *gin.Context) {
	passwd := c.PostForm("passwd")
	slug := c.Param("slug")

	content := content.GetContent(slug, "", "publish")
	fmt.Println(content)
	if content.Password != passwd {
		service.Error404(c)
		return
	}

	fmt.Println(c.Request.RequestURI)
	fmt.Println(c.Request.Referer())

	var path string
	if content.Type == "post" {
		path = "/post/" + content.Slug
	} else {
		path = "/page/" + content.Slug
	}

	c.SetCookie("protectPassword_"+content.Slug, passwd, 60*60*2, path, "", true, true)

	c.Redirect(http.StatusMovedPermanently, c.Request.Referer())
}
