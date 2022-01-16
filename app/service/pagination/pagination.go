/*
 * @Author: jia
 * @LastEditTime: 2021-11-12 21:30:38
 * @FilePath: /thinkecho/app/service/pagination/pagination.go
 * @Date: 2021-11-05 23:15:03
 * @Software: VS Code
 */
package pagination

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
)

type Pagination struct {
	TotalItems   int   // 所有条目数
	CurrentPage  int   // 当前页
	PerPageItems int   // 每页条目数
	pinFirstPage bool  // 显示第一页
	pinLastPage  bool  // 显示最后一页
	totalPages   int   // 所有分页总数
	pages        []int // 显示页码
}

func (p *Pagination) generatePagination() {
	//计算分页数
	p.totalPages = int(math.Ceil(float64(p.TotalItems) / float64(p.PerPageItems)))

	// 如果当前页面大于总分页数则使当前页面显示分页的最后一页
	if p.CurrentPage > p.totalPages {
		p.CurrentPage = p.totalPages
	}

	// 如果当前分页小于等于0，则使页面显示第一页
	if p.CurrentPage <= 0 {
		p.CurrentPage = 1
	}

	//获取第一页和最后一页页码
	firstPage := p.CurrentPage - 3
	lastPage := p.CurrentPage + 3

	//第一页页码不能小于1
	if firstPage < 1 {
		firstPage = 1
	}

	//最后一页页码不能大于总分页数
	if lastPage > p.totalPages {
		lastPage = p.totalPages
	}

	if p.totalPages > 6 {
		if lastPage < p.totalPages && p.CurrentPage <= 3 {
			lastPage = firstPage + 6 - 1
		}

		if p.CurrentPage > p.totalPages-6 {
			firstPage = lastPage - 6 + 1
		}
	} else {
		firstPage = 1
	}

	if firstPage != 1 {
		p.pinFirstPage = true
	}

	if lastPage != p.totalPages {
		p.pinLastPage = true
	}

	p.pages = make([]int, 0, lastPage-firstPage+1)
	for i := firstPage; i <= lastPage; i++ {
		p.pages = append(p.pages, i)
	}
}

/**
 * @description: 生成Html
 * @param {string} uri
 * @return {*}
 */
func (p *Pagination) HTML(uri string) template.HTML {
	p.generatePagination()

	var b bytes.Buffer
	if p.pinFirstPage {
		b.WriteString(`<a class="page-first" href="` + fmt.Sprintf("%s%d", uri, 1) + `">`)
		b.WriteString("1")
		b.WriteString(`</a>`)
		b.WriteString(`<span class="page-ellipsis">...</span>`)
	}

	for _, page := range p.pages {
		c := ""
		if p.CurrentPage == page {
			c = " current"
		}
		b.WriteString(`<a class="page-item` + c + `" href="` + fmt.Sprintf("%s%d", uri, page) + `">`)
		b.WriteString(fmt.Sprintf("%d", page))
		b.WriteString(`</a> `)
	}

	if p.pinLastPage {
		b.WriteString(`<span class="page-ellipsis">...</span> `)
		b.WriteString(`<a class="page-last" href="` + fmt.Sprintf("%s%d", uri, p.totalPages) + `">`)
		b.WriteString(fmt.Sprintf("%d", p.totalPages))
		b.WriteString(`</a> `)
	}

	return template.HTML(b.String())
}
