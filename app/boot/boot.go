/*
 * @Author: jia
 * @LastEditTime: 2022-01-16 22:57:44
 * @FilePath: /thinkecho/app/boot/boot.go
 * @Date: 2022-01-16 22:23:27
 * @Software: VS Code
 */
package boot

import (
	"fmt"
	"os"
	"thinkecho/app/model"
	"thinkecho/app/pkg/config"
)

func init() {
	var err error

	if err = config.InitConfig(); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if err = model.InitDb(); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}
