/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 23:05:26
 * @FilePath: /thinkecho/app/model/database/model.go
 * @Date: 2022-01-16 22:42:21
 * @Software: VS Code
 */
package database

import (
	"thinkecho/app/model"

	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	DB = model.DB
}
