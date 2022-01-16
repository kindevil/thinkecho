/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 23:30:16
 * @FilePath: /thinkecho/app/pkg/theme/func.go
 * @Date: 2021-11-10 14:14:25
 * @Software: VS Code
 */
package theme

import (
	"fmt"
	"html/template"
	"path"
	"thinkecho/app/database/content"
	"thinkecho/app/service/archive"
	"thinkecho/app/service/option"
)

var archiveTitleTemplate = map[string]string{
	"category": "分类 %s 下的文章",
	"tag":      "标签 %s 下的文章",
	"search":   "包含关键字 %s 的文章",
	"author":   "%s 发布的文章",
}

/**
 * @description: 页面标题
 * @param {...interface{}} separators
 * @return {*}
 */
func title() string {
	return option.GetOption("title")
}

/**
 * @description: 页面标题
 * @param {...interface{}} data
 * @return {*}
 */
func archiveTitle(data ...interface{}) string {
	var title = theme.ArchiveTitle

	if title == "" {
		return ""
	}

	switch theme.PageType {
	case "category":
		title = fmt.Sprintf(archiveTitleTemplate["category"], title)
	case "tag":
		title = fmt.Sprintf(archiveTitleTemplate["tag"], title)
	case "search":
		title = fmt.Sprintf(archiveTitleTemplate["search"], title)
	case "author":
		title = fmt.Sprintf(archiveTitleTemplate["author"], title)
	}

	if len(data) > 0 {
		title += data[0].(string)
	}

	return title
}

/**
 * @description: 网站描述
 * @param {*}
 * @return {*}
 */
func description() string {
	return option.GetOption("description")
}

/**
 * @description: 页面编码
 * @param {*}
 * @return {*}
 */
func charset() string {
	return option.GetOption("charset")
}

/**
 * @description: 网址网址
 * @param {*}
 * @return {*}
 */
func siteUrl(data ...interface{}) string {
	var uri string
	if len(data) > 0 {
		uri = data[0].(string)
	}
	return option.GetOption("siteUrl") + "/" + uri
}

/**
 * @description: logo链接
 * @param {*}
 * @return {*}
 */
func logoUrl() string {
	return option.GetOption("siteUrl")
}

/**
 * @description: 主题地址
 * @param {*}
 * @return {*}
 */
func themeUrl(data ...interface{}) string {
	var uri string
	if len(data) > 0 {
		uri = data[0].(string)
	}
	return siteUrl(path.Join(theme.themeName, "/", uri))
}

/**
 * @description: 头部信息
 * @param {*}
 * @return {*}
 */
func header() string {
	return ""
}

/**
 * @description: footer信息
 * @param {*}
 * @return {*}
 */
func footer() string {
	return ""
}

/**
 * @description: 判断当前页面
 * @param {string} page
 * @param {...interface{}} data
 * @return {*}
 */
func is(page string, data ...interface{}) bool {
	if page == theme.PageType {
		if len(data) > 0 {
			if data[0].(string) == theme.PageSlug {
				return true
			}
		} else {
			return true
		}
	}

	return false
}

/**
 * @description: 下一篇文章
 * @param {*}
 * @return {*}
 */
func next() template.HTML {
	if is("post") {
		data := content.GetContent(theme.PageSlug, "post", "publish")
		if !data.IsEmpty() && data.Cid != 0 {
			next := content.GetContentNext(data.Created)
			var nextLink string = "没有了"
			if !next.IsEmpty() && next.Cid != 0 {
				nextLink = fmt.Sprintf(`<a href="%s/post/%s">%s</a>`, option.GetOption("siteUrl"), next.Slug, next.Title)
			}
			return template.HTML(nextLink)
		}
	}

	return ""
}

/**
 * @description: 相关文章
 * @param {*}
 * @return {*}
 */
func related(args ...interface{}) []*archive.Post {
	if is("post") {
		data := content.GetContent(theme.PageSlug, "post", "publish")
		var mids []uint
		for _, v := range data.Tags {
			mids = append(mids, v.Mid)
		}

		var limit = 5

		if len(args) > 0 {
			limit = args[0].(int)
		}

		relateds := content.GetRelated(mids, data.Cid, limit)
		var posts []*archive.Post
		for _, post := range *relateds {
			posts = append(posts, archive.NewPost(theme.Context, &post))
		}

		return posts
	}
	return nil
}

/**
 * @description: 上一篇文章
 * @param {*}
 * @return {*}
 */
func prev() template.HTML {
	if is("post") {
		data := content.GetContent(theme.PageSlug, "post", "publish")
		if !data.IsEmpty() && data.Cid != 0 {
			prev := content.GetContentPrev(data.Created)
			var prevLink string = "没有了"
			if !prev.IsEmpty() && prev.Cid != 0 {
				prevLink = fmt.Sprintf(`<a href="%s/post/%s">%s</a>`, option.GetOption("siteUrl"), prev.Slug, prev.Title)
			}
			return template.HTML(prevLink)
		}
	}

	return ""
}

/**
 * @description: 页面导航
 * @param {*}
 * @return {*}
 */
type Page struct {
	Title     string
	Permalink string
}

/**
 * @description: 页面导航
 * @param {*}
 * @return {*}
 */
func pageList() []*Page {
	var pageList []*Page
	var pages = content.GetPageTitle()

	for _, page := range *pages {
		siteNav := &Page{
			Title:     page.Title,
			Permalink: siteUrl() + "page/" + page.Slug,
		}
		pageList = append(pageList, siteNav)
	}
	return pageList
}
