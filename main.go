/*
 * @Author: jia
 * @LastEditTime: 2021-11-26 18:08:55
 * @FilePath: /thinkecho/main.go
 * @Date: 2021-11-03 22:31:47
 * @Software: VS Code
 */
package main

import (
	"fmt"
	_ "thinkecho/app/boot"
	"thinkecho/app/database/content"
)

/**
 * @description: 主函数
 * @param {*}
 * @return {*}
 */
func main() {
	relateds := content.GetRelated([]uint{78}, 504, 1)
	fmt.Println(relateds)
	//server.Run()
}
