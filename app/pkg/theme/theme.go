/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 23:06:17
 * @FilePath: /thinkecho/app/pkg/theme/theme.go
 * @Date: 2021-11-09 21:46:39
 * @Software: VS Code
 */
package theme

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
	"thinkecho/app/service/option"

	_ "github.com/flosch/pongo2-addons"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Theme struct {
	themeFolder string
	themeName   string
	delims      struct {
		left  string
		right string
	}
	Template     *template.Template // 模版引擎
	FuncMap      map[string]interface{}
	ArchiveTitle string
	PageType     string
	PageSlug     string
	Context      *gin.Context
}

var (
	themeFolder     = "content/themes/" // 主题存放目录
	theme           *Theme
	htmlContentType = []string{"text/html; charset=utf-8"}
)

func init() {
	theme = &Theme{
		themeFolder: themeFolder,
		themeName:   option.GetOption("theme"),
		delims: struct {
			left  string
			right string
		}{"{{", "}}"},
		FuncMap: make(map[string]interface{}),
	}

	// 加载功能
	theme.LoadFunc()

	pattern := path.Join(theme.themeFolder, "/", theme.themeName, "/*.tmpl")
	theme.Template = template.Must(template.New("").Delims(theme.delims.left, theme.delims.right).Funcs(theme.FuncMap).ParseGlob(pattern))

	// 打印出读取到的模版文件
	if viper.GetString("mode") == "debug" {
		var buf strings.Builder
		for _, tmpl := range theme.Template.Templates() {
			buf.WriteString("\t- ")
			buf.WriteString(tmpl.Name())
			buf.WriteString("\n")
		}
		fmt.Printf("Loaded Theme Templates (%d): \n%s\n", len(theme.Template.Templates()), buf.String())
	}
}

/**
 * @description: 加载功能
 * @param {*}
 * @return {*}
 */
func (theme *Theme) LoadFunc() {
	theme.FuncMap["title"] = title
	theme.FuncMap["archiveTitle"] = archiveTitle
	theme.FuncMap["charset"] = charset
	theme.FuncMap["description"] = description
	theme.FuncMap["siteUrl"] = siteUrl
	theme.FuncMap["themeUrl"] = themeUrl
	theme.FuncMap["logoUrl"] = logoUrl
	theme.FuncMap["header"] = header
	theme.FuncMap["footer"] = footer
	theme.FuncMap["pageList"] = pageList
	theme.FuncMap["is"] = is
	theme.FuncMap["next"] = next
	theme.FuncMap["prev"] = prev
	theme.FuncMap["related"] = related
}

/**
 * @description: 添加功能
 * @param {string} name
 * @param {interface{}} fn
 * @return {*}
 */
func AddFunc(name string, fn interface{}) {
	theme.FuncMap[name] = fn
}

/**
 * @description: 返回当前主题目录
 * @param {*}
 * @return {*}
 */
func ThemeFolder() string {
	return path.Join(theme.themeFolder, "/", theme.themeName)
}

/**
 * @description: 加载模版
 * @param {string} name
 * @param {interface{}} data
 * @return {*}
 */
func (theme *Theme) Load(name string, data gin.H) *Html {
	html := &Html{
		Name: name,
		Data: data,
	}

	if viper.GetString("mode") == "debug" {
		html.Template = theme.loadTemplate()
	} else {
		html.Template = theme.Template
	}

	return html
}

//loadTemplate 动态加载模版
func (t *Theme) loadTemplate() *template.Template {
	pattern := path.Join(theme.themeFolder, "/", theme.themeName, "/*.tmpl")
	if pattern != "" {
		return template.Must(template.New("").Delims(t.delims.left, t.delims.right).Funcs(t.FuncMap).ParseGlob(pattern))
	}
	panic("the HTML debug render was created without files or glob pattern")
}

type Html struct {
	Template *template.Template
	Name     string
	Data     gin.H
}

/**
 * @description: 渲染模版
 * @param {http.ResponseWriter} w
 * @return {*}
 */
func (h *Html) Render(w http.ResponseWriter) error {
	h.WriteContentType(w)
	if h.Name == "" {
		return h.Template.Execute(w, h.Data)
	}
	return h.Template.ExecuteTemplate(w, h.Name, h.Data)
}

/**
 * @description: 设置http头
 * @param {http.ResponseWriter} w
 * @return {*}
 */
func (r *Html) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = htmlContentType
	}
}

/**
 * @description: 渲染页面
 * @param {*gin.Context} c
 * @param {int} code
 * @param {string} name
 * @param {interface{}} data
 * @return {*}
 */
func HTML(c *gin.Context, code int, name string, data gin.H) {
	theme.Context = c
	theme.ArchiveTitle = c.GetString("archiveTitle")
	theme.PageType = c.GetString("pageType")
	theme.PageSlug = c.GetString("pageSlug")

	html := theme.Load(name, data)

	c.Status(code)

	if !bodyAllowedForStatus(code) {
		html.WriteContentType(c.Writer)
		c.Writer.WriteHeaderNow()
		return
	}

	if err := html.Render(c.Writer); err != nil {
		panic(err)
	}
}

/**
 * @description: 判断页面状态
 * @param {int} status
 * @return {*}
 */
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
