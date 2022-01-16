/*
 * @Author: jia
 * @LastEditTime: 2021-11-04 21:58:05
 * @FilePath: /thinkecho/app/boot/boot.go
 * @Date: 2021-11-03 23:03:52
 * @Software: VS Code
 */
package boot

import (
	"fmt"
	"os"
	"thinkecho/app/database"
	"thinkecho/app/pkg/config"
)

func init() {
	var err error

	if err = config.InitConfig(); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if err = database.InitDb(); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}
