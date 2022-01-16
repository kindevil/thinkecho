/*
 * @Author: jia
 * @LastEditTime: 2021-11-15 21:26:32
 * @FilePath: /thinkecho/app/server/engine.go
 * @Date: 2021-11-10 15:38:55
 * @Software: VS Code
 */
package server

import (
	"thinkecho/app/frontend/author"
	"thinkecho/app/frontend/category"
	"thinkecho/app/frontend/home"
	"thinkecho/app/frontend/page"
	"thinkecho/app/frontend/password"
	"thinkecho/app/frontend/post"
	"thinkecho/app/frontend/search"
	"thinkecho/app/frontend/tag"
	"thinkecho/app/pkg/theme"
	"thinkecho/app/service"
	"thinkecho/app/service/option"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Server struct {
	Engine *gin.Engine
	Addr   string
}

var server *Server

func init() {
	if viper.GetString("mode") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server = &Server{
		Engine: gin.New(),
		Addr:   viper.GetString("address"),
	}
}

func (s *Server) route() {
	s.Engine.GET("/", home.Home)
	s.Engine.GET("/p/:num", home.Home)

	s.Engine.GET("/page/:slug", page.Page)
	s.Engine.GET("/post/:slug", post.Post)

	s.Engine.GET("/tag/:slug", tag.Tag)
	s.Engine.GET("/tag/:slug/p/:num", tag.Tag)

	s.Engine.GET("/category/:slug", category.Category)
	s.Engine.GET("/category/:slug/p/:num", category.Category)

	s.Engine.GET("/search/:keyword", search.Search)
	s.Engine.GET("/search/:keyword/p/:num", search.Search)

	s.Engine.GET("/author/:uid", author.Author)
	s.Engine.GET("/author/:uid/p/:num", author.Author)

	s.Engine.POST("/password/:slug", password.Password)

	s.Engine.NoRoute(service.Error404)
}

func Run() {
	server.Engine.Use(gin.Logger())
	server.Engine.Use(gin.Recovery())

	// 设置静态目录
	server.Engine.Static("/"+option.GetOption("theme"), theme.ThemeFolder())
	server.Engine.Static("/usr/uploads", "content/uploads")

	// 设置路由
	server.route()

	server.Engine.Run(server.Addr)
}
