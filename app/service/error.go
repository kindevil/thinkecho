/*
 * @Author: jia
 * @LastEditTime: 2021-11-10 16:00:11
 * @FilePath: /thinkecho/app/service/error.go
 * @Date: 2021-11-10 15:53:23
 * @Software: VS Code
 */
package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error404(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
	})
}
