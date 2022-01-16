/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:42:08
 * @FilePath: /thinkecho/app/frontend/category/category.go
 * @Date: 2021-11-15 13:41:21
 * @Software: VS Code
 */
package category

import (
	"net/http"
	"strconv"
	"thinkecho/app/database/content"
	"thinkecho/app/database/meta"
	"thinkecho/app/pkg/theme"
	"thinkecho/app/service"
	"thinkecho/app/service/archive"
	"thinkecho/app/service/option"
	"thinkecho/app/service/pagination"

	"github.com/gin-gonic/gin"
)

func Category(c *gin.Context) {
	var err error
	var currentPage int
	var pageSize int

	slug := c.Param("slug")

	category := meta.GetMeta(slug, "category")
	if category.IsEmpty() {
		service.Error404(c)
		return
	}

	c.Set("archiveTitle", category.Name)
	c.Set("pageType", "category")
	c.Set("pageSlug", category.Slug)

	// 获取当前页码
	currentPage, err = strconv.Atoi(c.Param("num"))
	if err != nil {
		currentPage = 1
	}

	pageSize = option.GetOptionInt("pageSize")

	//计算偏移量
	offset := (currentPage - 1) * pageSize

	//获取文章cid
	var cids []int
	for _, r := range category.Relationship {
		cids = append(cids, int(r.Cid))
	}

	//从数据库获取文章
	contents := content.GetContents("post", cids, "publish", 0, "", pageSize, offset, "desc")
	//获取总数量
	totalContents := content.GetContentCount("post", cids, "publish", 0, "")

	pageNav := &pagination.Pagination{
		TotalItems:   int(totalContents),
		CurrentPage:  currentPage,
		PerPageItems: pageSize,
	}

	var posts []*archive.Post
	for _, content := range *contents {
		posts = append(posts, archive.NewPost(c, &content))
	}

	uri := option.GetOption("siteUrl") + "/category/" + slug + "/p/"

	theme.HTML(c, http.StatusOK, "archive.tmpl", gin.H{
		"category": category,
		"archives": posts,
		"pageNav":  pageNav.HTML(uri),
	})
}
