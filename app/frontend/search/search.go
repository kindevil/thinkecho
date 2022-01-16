/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:41:43
 * @FilePath: /thinkecho/app/frontend/search/search.go
 * @Date: 2021-11-15 14:37:38
 * @Software: VS Code
 */
package search

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

func Search(c *gin.Context) {
	var err error
	var currentPage int
	var pageSize int

	keyword := c.Param("keyword")

	c.Set("archiveTitle", keyword)
	c.Set("pageType", "search")
	c.Set("pageSlug", keyword)

	// 获取当前页码
	currentPage, err = strconv.Atoi(c.Param("num"))
	if err != nil {
		currentPage = 1
	}

	pageSize = option.GetOptionInt("pageSize")

	//计算偏移量
	offset := (currentPage - 1) * pageSize

	//从数据库获取文章
	contents := content.GetContents("post", nil, "publish", 0, keyword, pageSize, offset, "desc")

	//获取总数量
	totalContents := content.GetContentCount("post", nil, "publish", 0, keyword)

	pageNav := &pagination.Pagination{
		TotalItems:   int(totalContents),
		CurrentPage:  currentPage,
		PerPageItems: pageSize,
	}

	var posts []*archive.Post
	for _, content := range *contents {
		posts = append(posts, archive.NewPost(c, &content))
	}

	uri := option.GetOption("siteUrl") + "/search/" + keyword + "/p/"

	theme.HTML(c, http.StatusOK, "archive.tmpl", gin.H{
		"archives": posts,
		"pageNav":  pageNav.HTML(uri),
	})
}
