/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:42:15
 * @FilePath: /thinkecho/app/frontend/author/author.go
 * @Date: 2021-11-15 17:32:08
 * @Software: VS Code
 */
package author

import (
	"fmt"
	"net/http"
	"strconv"
	"thinkecho/app/database/content"
	"thinkecho/app/database/user"
	"thinkecho/app/pkg/theme"
	"thinkecho/app/service"
	"thinkecho/app/service/archive"
	"thinkecho/app/service/option"
	"thinkecho/app/service/pagination"

	"github.com/gin-gonic/gin"
)

func Author(c *gin.Context) {
	var err error
	var currentPage int
	var pageSize int

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		uid = 0
	}

	user := user.GetUserByID(uid)
	if user.IsEmpty() || user.UID == 0 {
		service.Error404(c)
		return
	}

	fmt.Println(user)

	c.Set("archiveTitle", user.ScreenName)
	c.Set("pageType", "author")
	c.Set("pageSlug", user.UID)

	// 获取当前页码
	currentPage, err = strconv.Atoi(c.Param("num"))
	if err != nil {
		currentPage = 1
	}

	pageSize = option.GetOptionInt("pageSize")

	//计算偏移量
	offset := (currentPage - 1) * pageSize

	//从数据库获取文章
	contents := content.GetContents("post", nil, "publish", user.UID, "", pageSize, offset, "desc")
	//获取总数量
	totalContents := content.GetContentCount("post", nil, "publish", user.UID, "")

	pageNav := &pagination.Pagination{
		TotalItems:   int(totalContents),
		CurrentPage:  currentPage,
		PerPageItems: pageSize,
	}

	var posts []*archive.Post
	for _, content := range *contents {
		posts = append(posts, archive.NewPost(c, &content))
	}

	uri := option.GetOption("siteUrl") + fmt.Sprintf("/author/%d/p/", uid)

	theme.HTML(c, http.StatusOK, "archive.tmpl", gin.H{
		"author":   user,
		"archives": posts,
		"pageNav":  pageNav.HTML(uri),
	})
}
