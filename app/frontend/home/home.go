/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:42:03
 * @FilePath: /thinkecho/app/frontend/home/home.go
 * @Date: 2021-11-12 17:11:52
 * @Software: VS Code
 */
package home

import (
	"net/http"
	"strconv"
	"thinkecho/app/database/content"
	"thinkecho/app/pkg/theme"
	"thinkecho/app/service/archive"
	"thinkecho/app/service/option"
	"thinkecho/app/service/pagination"

	"github.com/gin-gonic/gin"
)

func Archives(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "200"})
}

func Home(c *gin.Context) {
	var err error
	var currentPage int
	var pageSize int

	c.Set("archiveTitle", "")
	c.Set("pageType", "index")
	c.Set("pageSlug", "")

	// 获取当前页码
	currentPage, err = strconv.Atoi(c.Param("num"))
	if err != nil {
		currentPage = 1
	}

	pageSize = option.GetOptionInt("pageSize")

	//计算偏移量
	offset := (currentPage - 1) * pageSize

	//从数据库获取文章
	contents := content.GetContents("post", nil, "publish", 0, "", pageSize, offset, "desc")

	//获取总数量
	totalContents := content.GetContentCount("post", nil, "publish", 0, "")

	pageNav := &pagination.Pagination{
		TotalItems:   int(totalContents),
		CurrentPage:  currentPage,
		PerPageItems: pageSize,
	}

	var posts []*archive.Post
	for _, content := range *contents {
		posts = append(posts, archive.NewPost(c, &content))
	}

	uri := option.GetOption("siteUrl") + "/p/"

	theme.HTML(c, http.StatusOK, "index.tmpl", gin.H{
		"archives": posts,
		"pageNav":  pageNav.HTML(uri),
	})
}
