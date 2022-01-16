/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 22:30:16
 * @FilePath: /thinkecho/app/pkg/file/file.go
 * @Date: 2022-01-16 22:30:16
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
