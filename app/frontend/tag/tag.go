/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:41:52
 * @FilePath: /thinkecho/app/frontend/tag/tag.go
 * @Date: 2021-11-15 10:34:27
 * @Software: VS Code
 */
package tag

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

func Tag(c *gin.Context) {
	var err error
	var currentPage int
	var pageSize int

	slug := c.Param("slug")

	tag := meta.GetMeta(slug, "tag")
	if tag.IsEmpty() {
		service.Error404(c)
		return
	}

	c.Set("archiveTitle", tag.Name)
	c.Set("pageType", "tag")
	c.Set("pageSlug", tag.Slug)

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
	for _, r := range tag.Relationship {
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

	uri := option.GetOption("siteUrl") + "/tag/" + slug + "/p/"

	theme.HTML(c, http.StatusOK, "archive.tmpl", gin.H{
		"tag":      tag,
		"archives": posts,
		"pageNav":  pageNav.HTML(uri),
	})
}
