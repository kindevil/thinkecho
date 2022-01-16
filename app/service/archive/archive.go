/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:45:13
 * @FilePath: /thinkecho/app/service/archive/archive.go
 * @Date: 2021-11-06 10:47:49
 * @Software: VS Code
 */
package archive

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"thinkecho/app/database/field"
	"thinkecho/app/database/meta"
	"thinkecho/app/service/option"
	"time"

	"github.com/88250/lute"
	"github.com/gin-gonic/gin"
)

type Archive struct {
	Cid          uint
	Title        string
	Slug         string
	Author       Author
	allowComment uint
	commentsNum  uint
	text         string
	password     string
	created      int64
	fields       []*field.Field
	context      *gin.Context
	isPost       bool
}

/**
 * @description: 链接
 * @param {*}
 * @return {*}
 */
func (a *Archive) Permalink() string {
	var b bytes.Buffer
	b.WriteString(option.GetOption("siteUrl"))

	if a.isPost {
		b.WriteString("/post/")
	} else {
		b.WriteString("/page/")
	}

	b.WriteString(a.Slug)
	return b.String()
}

/**
 * @description: 发布日期
 * @param {*}
 * @return {*}
 */
func (a *Archive) Date() string {
	nowYear := time.Now().Year()
	year := time.Unix(a.created, 0).Year()

	if nowYear-year == 0 {
		return time.Unix(a.created, 0).Format("01月02日")
	}

	return time.Unix(a.created, 0).Format("2006年01月02日")
}

/**
 * @description: 内容输出
 * @param {string} readMore
 * @return {*}
 */
func (a *Archive) Content(more ...interface{}) template.HTML {
	var readmore string
	if len(more) > 0 {
		readmore = more[0].(string)
	}

	passwd, err := a.context.Cookie("protectPassword_" + a.Slug)
	if err != nil {
		passwd = ""
	}

	if a.password == "" || a.password == passwd {
		return a.markdown(readmore)
	}

	return a.passwordForm()
}

/**
 * @description: markdown格式文档输出
 * @param {string} readMore
 * @return {*}
 */
func (a *Archive) markdown(more string) template.HTML {
	luteEngine := lute.New()
	html := luteEngine.MarkdownStr(a.Title, a.text)
	if more == "" {
		return template.HTML(html)
	} else {
		return template.HTML(strings.Split(html, "<!--more-->")[0] + "<p class=\"more\"><a href=\"" + a.Permalink() + "\" title=\"" + a.Title + "\">" + more + "</a></p>")
	}
}

/**
 * @description: 密码输入对话框
 * @param {*}
 * @return {*}
 */
func (a *Archive) passwordForm() template.HTML {
	var b bytes.Buffer
	b.WriteString(`<form action="` + option.GetOption("siteUrl") + `/password/` + a.Slug + `" class="password" method="post" enctype="application/x-www-form-urlencoded">`)
	b.WriteString(`<input type="password" name="passwd" required="required"/>`)
	b.WriteString(`<button type="submit">提交</button>`)
	b.WriteString(`</form>`)
	return template.HTML(b.String())
}

/**
 * @description: 评论链接
 * @param {*}
 * @return {*}
 */
func (a *Archive) CommentUrl() string {
	return a.Permalink() + "#comments"
}

/**
 * @description: 评论数输出
 * @param {*}
 * @return {*}
 */
func (a *Archive) CommentsNum() template.HTML {
	var html string
	if a.commentsNum == 0 {
		html = "无评论"
	}

	if a.commentsNum == 1 {
		html = "1条评论"
	}

	if a.commentsNum > 1 {
		html = fmt.Sprintf("%d条评论", a.commentsNum)
	}
	return template.HTML(html)
}

/**
 * @description: 判断是否允许评论
 * @param {*}
 * @return {*}
 */
func (a *Archive) AllowComment() bool {
	return a.allowComment == 1
}

/**
 * @description: 获取字符串字段值
 * @param {string} name
 * @return {*}
 */
func (a *Archive) GetFieldStr(name string) string {
	for _, v := range a.fields {
		if v.Name == name && v.Type == "str" {
			return v.StrValue
		}
	}
	return ""
}

/**
 * @description: 获取整数字段值
 * @param {string} name
 * @return {*}
 */
func (a *Archive) GetFieldInt(name string) int {
	for _, v := range a.fields {
		if v.Name == name && v.Type == "int" {
			return v.IntValue
		}
	}
	return 0
}

/**
 * @description:获取浮点型字段值
 * @param {string} name
 * @return {*}
 */
func (a *Archive) GetFieldFloat(name string) float64 {
	for _, v := range a.fields {
		if v.Name == name && v.Type == "float" {
			return v.FloatValue
		}
	}
	return 0
}

type Author struct {
	Name string
	uid  uint
}

/**
 * @description: 生成用户链接
 * @param {*}
 * @return {*}
 */
func (a *Author) Permalink() string {
	return option.GetOption("siteUrl") + fmt.Sprintf("/author/%d", a.uid)
}

type Meta struct {
	categories []*meta.Meta
	tags       []*meta.Meta
}

/**
 * @description: 输出分类
 * @param {string} separator
 * @return {*}
 */
func (m *Meta) Category(data ...interface{}) template.HTML {
	var separator string
	if len(data) > 0 {
		separator = data[0].(string)
	}
	return m.meta(m.categories, separator, false)
}

/**
 * @description: 输出标签
 * @param {string} separator
 * @return {*}
 */
func (m *Meta) Tags(data ...interface{}) template.HTML {
	var separator string
	if len(data) > 0 {
		separator = data[0].(string)
	}
	return m.meta(m.tags, separator, true)
}

/**
 * @description: 输出分类
 * @param {string} separator
 * @return {*}
 */
func (m *Meta) meta(metas []*meta.Meta, separator string, isTag bool) template.HTML {
	if m.categories == nil {
		return ""
	}

	var defaultSeparator = ","
	if separator != "" {
		defaultSeparator = separator
	}

	var metaSlice []string
	for _, meta := range metas {
		var b bytes.Buffer
		b.WriteString(`<a href=`)
		b.WriteString(option.GetOption("siteUrl"))
		if isTag {
			b.WriteString(`/tag/`)
		} else {
			b.WriteString(`/category/`)
		}
		b.WriteString(meta.Slug)
		b.WriteString(`>` + meta.Name + `</a>`)

		metaSlice = append(metaSlice, b.String())
	}

	return template.HTML(strings.Join(metaSlice, defaultSeparator))
}
