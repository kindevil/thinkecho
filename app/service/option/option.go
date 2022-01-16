/*
 * @Author: jia
 * @LastEditTime: 2021-11-09 21:15:58
 * @FilePath: /thinkecho/app/service/option/option.go
 * @Date: 2021-11-09 21:09:41
 * @Software: VS Code
 */
package option

import (
	"strconv"
	"thinkecho/app/database/option"
)

var Options map[string]string

func init() {
	Options = make(map[string]string)
	options := option.GetOptions()
	for _, option := range *options {
		Options[option.Name] = option.Value
	}
}

/**
 * @description: 获取option
 * @param {string} name
 * @return {string}
 */
func GetOption(name string) string {
	if value, ok := Options[name]; ok {
		return value
	}
	return ""
}

/**
 * @description: 获取option,返回值为int
 * @param {string} name
 * @return {int}
 */
func GetOptionInt(name string) int {
	value, err := strconv.Atoi(GetOption(name))
	if err != nil {
		return 0
	}

	return value
}

/**
 * @description: 获取option,返回值为string
 * @param {string} name
 * @return {string}
 */
func GetOptionString(name string) string {
	return GetOption(name)
}
