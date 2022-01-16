/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 19:12:07
 * @FilePath: /thinkecho/main.go
 * @Date: 2021-11-03 22:31:47
 * @Software: VS Code
 */
package main

import (
	_ "thinkecho/app/boot"
	"thinkecho/app/server"
)

/**
 * @description: 主函数
 * @param {*}
 * @return {*}
 */
func main() {
	server.Run()
}
