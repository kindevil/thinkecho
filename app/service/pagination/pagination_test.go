/*
 * @Author: jia
 * @LastEditTime: 2021-11-06 10:46:42
 * @FilePath: /thinkecho/app/service/pagination/pagination_test.go
 * @Date: 2021-11-05 23:52:50
 * @Software: VS Code
 */
package pagination_test

import (
	"fmt"
	"testing"
	"thinkecho/app/service/pagination"
)

func TestPagination(t *testing.T) {
	pagination := &pagination.Pagination{
		TotalItems:   48,
		CurrentPage:  2,
		PerPageItems: 6,
	}

	html := pagination.HTML("http://kindevil.com/page/")

	fmt.Println(html)
}
