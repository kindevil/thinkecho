/*
 * @Author: jia
 * @LastEditTime: 2021-11-03 22:52:11
 * @FilePath: /thinkecho/app/pkg/file/file.go
 * @Date: 2021-11-03 22:51:05
 * @Software: VS Code
 */
package file

import "os"

/**
 * @description: 目录或者文件是否存在
 * @param {string} path
 * @return {bool}
 */
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	return false
}
