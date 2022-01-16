// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import (
	"bytes"
	"github.com/88250/lute/lex"
	"strconv"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

func (lute *Lute) SpinBlockDOM(ivHTML string) (ovHTML string) {
	//fmt.Println(ivHTML)
	markdown := lute.blockDOM2Md(ivHTML)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)

	firstChild := tree.Root.FirstChild
	if ast.NodeParagraph == firstChild.Type && "" == firstChild.ID {
		if second := firstChild.Next; nil != second && nil != second.Next && ast.NodeKramdownBlockIAL == second.Next.Type {
			// 软换行后生成两个块，需要把老 ID 调整到第一个块上
			ial := second.Next
			firstChild.ID, second.ID = second.ID, ""
			firstChild.InsertAfter(ial)
			firstChild.KramdownIAL, second.KramdownIAL = second.KramdownIAL, nil
		}
	}
	if ast.NodeKramdownBlockIAL == firstChild.Type && nil != firstChild.Next && ast.NodeKramdownBlockIAL == firstChild.Next.Type && util.IsDocIAL(firstChild.Next.Tokens) {
		// 空段落块还原
		ialArray := parse.Tokens2IAL(firstChild.Tokens)
		ial := parse.IAL2Map(ialArray)
		p := &ast.Node{Type: ast.NodeParagraph, ID: ial["id"], KramdownIAL: ialArray}
		firstChild.InsertBefore(p)
	}

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) HTML2BlockDOM(sHTML string) (vHTML string) {
	//fmt.Println(sHTML)
	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.HTML2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	return
}

func (lute *Lute) BlockDOM2HTML(vhtml string) (sHTML string) {
	markdown := lute.blockDOM2Md(vhtml)
	sHTML = lute.Md2HTML(markdown)
	return
}

func (lute *Lute) Md2BlockDOM(markdown string) (vHTML string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	return
}

func (lute *Lute) InlineMd2BlockDOM(markdown string) (vHTML string) {
	tree := parse.Inline("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	return
}

func (lute *Lute) BlockDOM2Md(htmlStr string) (markdown string) {
	//fmt.Println(htmlStr)
	markdown = lute.blockDOM2Md(htmlStr)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

func (lute *Lute) BlockDOM2StdMd(htmlStr string) (markdown string) {
	htmlStr = strings.ReplaceAll(htmlStr, parse.Zwsp, "")

	// DOM 转 AST
	tree := lute.BlockDOM2Tree(htmlStr)

	// 将 kramdown IAL 节点内容置空
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		if ast.NodeKramdownBlockIAL == n.Type || ast.NodeKramdownSpanIAL == n.Type {
			n.Tokens = nil
		}
		return ast.WalkContinue
	})

	// 将 AST 进行 Markdown 格式化渲染
	options := render.NewOptions()
	options.AutoSpace = false
	options.FixTermTypo = false
	options.KramdownBlockIAL = true
	options.KramdownSpanIAL = true
	renderer := render.NewFormatRenderer(tree, options)
	formatted := renderer.Render()
	markdown = util.BytesToStr(formatted)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

func (lute *Lute) BlockDOM2Text(htmlStr string) (text string) {
	tree := lute.BlockDOM2Tree(htmlStr)
	return tree.Root.Text()
}

func (lute *Lute) BlockDOM2TextLen(htmlStr string) int {
	tree := lute.BlockDOM2Tree(htmlStr)
	return tree.Root.TextLen()
}

func (lute *Lute) Tree2BlockDOM(tree *parse.Tree, options *render.Options) (vHTML string) {
	renderer := render.NewBlockRenderer(tree, options)
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	vHTML = strings.ReplaceAll(vHTML, util.Caret, "<wbr>")
	return
}

func RenderNodeBlockDOM(node *ast.Node, parseOptions *parse.Options, renderOptions *render.Options) string {
	root := &ast.Node{Type: ast.NodeDocument}
	tree := &parse.Tree{Root: root, Context: &parse.Context{ParseOption: parseOptions}}
	renderer := render.NewBlockRenderer(tree, renderOptions)
	renderer.Writer = &bytes.Buffer{}
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := renderer.RendererFuncs[n.Type]
		return rendererFunc(n, entering)
	})
	return renderer.Writer.String()
}

func (lute *Lute) BlockDOM2Tree(htmlStr string) (ret *parse.Tree) {
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</strong>", "</strong>\n<wbr>\n")
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</em>", "</em>\n<wbr>\n")
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</s>", "</s>\n<wbr>\n")
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</u>", "</u>\n<wbr>\n")
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</span>", "</span>\n<wbr>\n")
	htmlStr = strings.ReplaceAll(htmlStr, "<wbr>", util.Caret)

	// 替换结尾空白，否则 HTML 解析会产生冗余节点导致生成空的代码块
	htmlStr = strings.ReplaceAll(htmlStr, "\t\n", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "    \n", "  \n")

	// 将字符串解析为 DOM 树
	htmlRoot := lute.parseHTML(htmlStr)
	if nil == htmlRoot {
		return
	}

	// 调整 DOM 结构
	lute.adjustVditorDOM(htmlRoot)

	// 将 HTML 树转换为 Markdown AST
	ret = &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument}, Context: &parse.Context{ParseOption: lute.ParseOptions}}
	ret.Context.Tip = ret.Root
	for c := htmlRoot.FirstChild; nil != c; c = c.NextSibling {
		lute.genASTByBlockDOM(c, ret)
	}

	// 调整树结构
	ast.Walk(ret.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeInlineHTML, ast.NodeCodeSpan, ast.NodeInlineMath, ast.NodeHTMLBlock, ast.NodeCodeBlockCode, ast.NodeMathBlockContent:
				n.Tokens = html.UnescapeHTML(n.Tokens)
				if nil != n.Next && ast.NodeCodeSpan == n.Next.Type && n.CodeMarkerLen == n.Next.CodeMarkerLen && nil != n.FirstChild && nil != n.FirstChild.Next {
					// 合并代码节点 https://github.com/Vanessa219/vditor/issues/167
					n.FirstChild.Next.Tokens = append(n.FirstChild.Next.Tokens, n.Next.FirstChild.Next.Tokens...)
					n.Next.Unlink()
				}
			case ast.NodeStrong, ast.NodeEmphasis, ast.NodeStrikethrough, ast.NodeUnderline:
				lute.mergeSameSpan(n, n.Type)
			}
		}
		return ast.WalkContinue
	})
	return
}

func (lute *Lute) mergeSameSpan(n *ast.Node, typ ast.NodeType) {
	if nil != n.Next && typ == n.Next.Type && nil != n.Next.Next && ast.NodeKramdownSpanIAL != n.Next.Next.Type {
		var spanChildren []*ast.Node
		n.Next.FirstChild.Unlink() // open marker
		n.Next.LastChild.Unlink()  // close marker
		for c := n.Next.FirstChild; nil != c; c = c.Next {
			spanChildren = append(spanChildren, c)
		}
		for _, c := range spanChildren {
			n.LastChild.InsertBefore(c)
		}
		n.Next.Unlink()
	}
}

func (lute *Lute) CancelSuperBlock(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	if ast.NodeSuperBlock != tree.Root.FirstChild.Type {
		return ivHTML
	}

	sb := tree.Root.FirstChild

	var blocks []*ast.Node
	for b := sb.FirstChild; nil != b; b = b.Next {
		blocks = append(blocks, b)
	}
	for _, b := range blocks {
		tree.Root.AppendChild(b)
	}
	sb.Unlink()

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) BlocksMergeSuperBlock(ivHTML string, layout string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)

	sb := &ast.Node{ID: ast.NewNodeID(), Type: ast.NodeSuperBlock}
	openMarkder := &ast.Node{Type: ast.NodeSuperBlockOpenMarker}
	sb.AppendChild(openMarkder)
	sb.AppendChild(&ast.Node{Type: ast.NodeSuperBlockLayoutMarker, Tokens: util.StrToBytes(layout)})
	var blocks []*ast.Node
	for n := tree.Root.FirstChild; nil != n; n = n.Next {
		blocks = append(blocks, n)
	}
	for _, b := range blocks {
		sb.AppendChild(b)
	}
	closeMarker := &ast.Node{Type: ast.NodeSuperBlockCloseMarker}
	sb.AppendChild(closeMarker)
	tree.Root.AppendChild(sb)

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) CancelList(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	if ast.NodeList != tree.Root.FirstChild.Type {
		return ivHTML
	}

	list := tree.Root.FirstChild

	var appends, unlinks []*ast.Node
	for li := list.FirstChild; nil != li; li = li.Next {
		for c := li.FirstChild; nil != c; c = c.Next {
			if ast.NodeTaskListItemMarker != c.Type {
				appends = append(appends, c)
			}
		}
		unlinks = append(unlinks, li)
	}
	for _, c := range appends {
		tree.Root.AppendChild(c)
	}
	for _, n := range unlinks {
		n.Unlink()
	}
	list.Unlink()

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) HLevel(ivHTML string, level string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	node := tree.Root.FirstChild
	if ast.NodeHeading != node.Type {
		return ivHTML
	}

	node.HeadingLevel, _ = strconv.Atoi(level)
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) H2P(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	node := tree.Root.FirstChild
	if ast.NodeHeading != node.Type {
		return ivHTML
	}

	node.Type = ast.NodeParagraph
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) P2H(ivHTML, level string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	node := tree.Root.FirstChild
	if ast.NodeParagraph != node.Type {
		return ivHTML
	}

	node.Type = ast.NodeHeading
	node.HeadingLevel, _ = strconv.Atoi(level)
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) Blocks2TLs(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)

	list := &ast.Node{ID: ast.NewNodeID(), Type: ast.NodeList, ListData: &ast.ListData{Typ: 3}}
	var lis, blocks, ials []*ast.Node
	for n := tree.Root.FirstChild; nil != n; n = n.Next {
		if ast.NodeKramdownBlockIAL != n.Type {
			li := &ast.Node{ID: ast.NewNodeID(), Type: ast.NodeListItem, ListData: &ast.ListData{Marker: []byte("*"), Typ: 3}}
			lis = append(lis, li)
			blocks = append(blocks, n)
		} else {
			ials = append(ials, n)
		}
	}
	for i, li := range lis {
		li.AppendChild(&ast.Node{Type: ast.NodeTaskListItemMarker})
		li.AppendChild(blocks[i])
		li.AppendChild(ials[i])
		liID := ast.NewNodeID()
		liIAL := &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: parse.IAL2Tokens([][]string{{"id", liID}})}
		li.SetIALAttr("id", liID)
		li.InsertAfter(liIAL)
		list.AppendChild(li)
	}
	tree.Root.AppendChild(list)

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) Blocks2OLs(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)

	list := &ast.Node{ID: ast.NewNodeID(), Type: ast.NodeList, ListData: &ast.ListData{Typ: 1, Start: 1}}
	var lis, blocks, ials []*ast.Node
	num := 1
	for n := tree.Root.FirstChild; nil != n; n = n.Next {
		if ast.NodeKramdownBlockIAL != n.Type {
			li := &ast.Node{ID: ast.NewNodeID(), Type: ast.NodeListItem, ListData: &ast.ListData{Marker: []byte("*"), Num: num, Typ: 1}}
			lis = append(lis, li)
			blocks = append(blocks, n)
			num++
		} else {
			ials = append(ials, n)
		}
	}
	for i, li := range lis {
		li.AppendChild(blocks[i])
		li.AppendChild(ials[i])
		liID := ast.NewNodeID()
		liIAL := &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: parse.IAL2Tokens([][]string{{"id", liID}})}
		li.SetIALAttr("id", liID)
		li.InsertAfter(liIAL)
		list.AppendChild(li)
	}
	tree.Root.AppendChild(list)

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) Blocks2ULs(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)

	list := &ast.Node{ID: ast.NewNodeID(), Type: ast.NodeList, ListData: &ast.ListData{}}
	var lis, blocks, ials []*ast.Node
	for n := tree.Root.FirstChild; nil != n; n = n.Next {
		if ast.NodeKramdownBlockIAL != n.Type {
			li := &ast.Node{ID: ast.NewNodeID(), Type: ast.NodeListItem, ListData: &ast.ListData{Marker: []byte("*")}}
			lis = append(lis, li)
			blocks = append(blocks, n)
		} else {
			ials = append(ials, n)
		}
	}
	for i, li := range lis {
		li.AppendChild(blocks[i])
		li.AppendChild(ials[i])
		liID := ast.NewNodeID()
		liIAL := &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: parse.IAL2Tokens([][]string{{"id", liID}})}
		li.SetIALAttr("id", liID)
		li.InsertAfter(liIAL)
		list.AppendChild(li)
	}
	tree.Root.AppendChild(list)

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) Blocks2Blockquote(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)

	var blocks, ials []*ast.Node
	for n := tree.Root.FirstChild; nil != n; n = n.Next {
		if ast.NodeKramdownBlockIAL != n.Type {
			blocks = append(blocks, n)
		} else {
			ials = append(ials, n)
		}
	}

	bqID := ast.NewNodeID()
	bq := &ast.Node{ID: bqID, Type: ast.NodeBlockquote}
	for i, _ := range blocks {
		bq.AppendChild(blocks[i])
		bq.AppendChild(ials[i])
	}
	bqIAL := &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: parse.IAL2Tokens([][]string{{"id", bqID}})}
	bq.InsertAfter(bqIAL)
	tree.Root.AppendChild(bq)
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) Blocks2Ps(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	node := tree.Root.FirstChild

	for p := node; nil != p; p = p.Next {
		if ast.NodeHeading == p.Type {
			p.Type = ast.NodeParagraph
		}
	}
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) Blocks2Hs(ivHTML, level string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	node := tree.Root.FirstChild

	for p := node; nil != p; p = p.Next {
		if ast.NodeParagraph == p.Type || ast.NodeHeading == p.Type {
			p.Type = ast.NodeHeading
			p.HeadingLevel, _ = strconv.Atoi(level)
		}
	}
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) OL2TL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)

	tree.Root.FirstChild.ListData.Typ = 3
	for li := tree.Root.FirstChild.FirstChild; nil != li; li = li.Next {
		if ast.NodeListItem == li.Type {
			li.ListData.Typ = 3
			li.PrependChild(&ast.Node{Type: ast.NodeTaskListItemMarker})
		}
	}
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) UL2TL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)

	tree.Root.FirstChild.ListData.Typ = 3
	for li := tree.Root.FirstChild.FirstChild; nil != li; li = li.Next {
		if ast.NodeListItem == li.Type {
			li.ListData.Typ = 3
			li.PrependChild(&ast.Node{Type: ast.NodeTaskListItemMarker})
		}
	}
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) TL2OL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	list := tree.Root.FirstChild
	if ast.NodeList != list.Type || 3 != list.ListData.Typ {
		return ivHTML
	}

	list.ListData.Typ = 1
	var unlinks []*ast.Node
	for li := list.FirstChild; nil != li; li = li.Next {
		if ast.NodeKramdownBlockIAL == li.Type {
			continue
		}
		unlinks = append(unlinks, li.FirstChild) // task marker
		li.ListData.Typ = 1
		if nil == li.Previous {
			li.ListData.Start = 1
			li.ListData.Num = 1
		}
	}
	for _, n := range unlinks {
		n.Unlink()
	}

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) TL2UL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	list := tree.Root.FirstChild
	if ast.NodeList != list.Type || 3 != list.ListData.Typ {
		return ivHTML
	}

	list.ListData.Typ = 0
	var unlinks []*ast.Node
	for li := list.FirstChild; nil != li; li = li.Next {
		if ast.NodeKramdownBlockIAL == li.Type {
			continue
		}
		unlinks = append(unlinks, li.FirstChild) // task marker
		li.ListData.Typ = 0
	}
	for _, n := range unlinks {
		n.Unlink()
	}

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) OL2UL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	list := tree.Root.FirstChild
	if ast.NodeList != list.Type {
		return ivHTML
	}

	list.ListData.Typ = 0
	for li := list.FirstChild; nil != li; li = li.Next {
		if ast.NodeKramdownBlockIAL == li.Type {
			continue
		}
		li.ListData.Typ = 0
	}

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) UL2OL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	list := tree.Root.FirstChild
	if ast.NodeList != list.Type {
		return ivHTML
	}

	num := 1
	list.ListData.Typ = 1
	for li := list.FirstChild; nil != li; li = li.Next {
		if ast.NodeKramdownBlockIAL == li.Type {
			continue
		}
		li.ListData.Typ = 1
		li.ListData.Num = num
		num++
	}

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) blockDOM2Md(htmlStr string) (markdown string) {
	tree := lute.BlockDOM2Tree(htmlStr)

	// 将 AST 进行 Markdown 格式化渲染
	options := render.NewOptions()
	options.AutoSpace = false
	options.FixTermTypo = false
	options.KramdownBlockIAL = true
	options.KramdownSpanIAL = true
	renderer := render.NewFormatRenderer(tree, options)
	formatted := renderer.Render()
	markdown = string(formatted)
	return
}

func (lute *Lute) genASTByBlockDOM(n *html.Node, tree *parse.Tree) {
	class := lute.domAttrValue(n, "class")
	if "protyle-attr" == class ||
		strings.Contains(class, "__copy") ||
		strings.Contains(class, "protyle-linenumber__rows") {
		return
	}

	if "1" == lute.domAttrValue(n, "spin") {
		return
	}

	if strings.Contains(class, "protyle-action") {
		if ast.NodeCodeBlock == tree.Context.Tip.Type {
			languageNode := n.FirstChild
			language := ""
			if nil != languageNode.FirstChild {
				language = languageNode.FirstChild.Data
			}
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker, CodeBlockInfo: util.StrToBytes(language)})
			code := lute.domText(n.NextSibling)
			lines := strings.Split(code, "\n")
			buf := bytes.Buffer{}
			for i, line := range lines {
				if strings.Contains(line, "```") {
					line = strings.ReplaceAll(line, "```", parse.Zwj+"```")
				} else {
					line = strings.ReplaceAll(line, parse.Zwj, "")
				}
				buf.WriteString(line)
				if i < len(lines)-1 {
					buf.WriteByte('\n')
				}
			}
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeCodeBlockCode, Tokens: buf.Bytes()})
		} else if ast.NodeListItem == tree.Context.Tip.Type {
			if 3 == tree.Context.Tip.ListData.Typ { // 任务列表
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeTaskListItemMarker, TaskListItemChecked: strings.Contains(lute.domAttrValue(n.Parent, "class"), "protyle-task--done")})
			}
		}
		return
	}

	if "true" == lute.domAttrValue(n, "contenteditable") {
		lute.genASTContenteditable(n, tree)
		return
	}

	dataType := ast.Str2NodeType(lute.domAttrValue(n, "data-type"))

	nodeID := lute.domAttrValue(n, "data-node-id")
	node := &ast.Node{ID: nodeID}
	if "" != node.ID && !lute.parentIs(n, atom.Table) {
		node.KramdownIAL = [][]string{{"id", node.ID}}
		ialTokens := lute.setBlockIAL(n, node)
		ial := &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: ialTokens}
		defer tree.Context.TipAppendChild(ial)
	}

	switch dataType {
	case ast.NodeBlockQueryEmbed:
		node.Type = ast.NodeBlockQueryEmbed
		node.AppendChild(&ast.Node{Type: ast.NodeOpenBrace})
		node.AppendChild(&ast.Node{Type: ast.NodeOpenBrace})
		content := lute.domAttrValue(n, "data-content")
		node.AppendChild(&ast.Node{Type: ast.NodeBlockQueryEmbedScript, Tokens: util.StrToBytes(content)})
		node.AppendChild(&ast.Node{Type: ast.NodeCloseBrace})
		node.AppendChild(&ast.Node{Type: ast.NodeCloseBrace})
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeTable:
		node.Type = ast.NodeTable
		var tableAligns []int
		if nil == n.FirstChild {
			return
		}

		tableDiv := n.FirstChild
		table := lute.domChild(tableDiv, atom.Table)
		if nil == table {
			return
		}
		for th := table.FirstChild.FirstChild.FirstChild; nil != th; th = th.NextSibling {
			align := lute.domAttrValue(th, "align")
			switch align {
			case "left":
				tableAligns = append(tableAligns, 1)
			case "center":
				tableAligns = append(tableAligns, 2)
			case "right":
				tableAligns = append(tableAligns, 3)
			default:
				tableAligns = append(tableAligns, 0)
			}
		}
		node.TableAligns = tableAligns
		node.Tokens = nil
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()

		lute.genASTContenteditable(table, tree)
		return
	case ast.NodeParagraph:
		node.Type = ast.NodeParagraph
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeHeading:
		text := lute.domText(n)
		if lute.parentIs(n, atom.Table) {
			node.Tokens = []byte(strings.TrimSpace(text))
			for bytes.HasPrefix(node.Tokens, []byte("#")) {
				node.Tokens = bytes.TrimPrefix(node.Tokens, []byte("#"))
			}
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeHeading
		level := lute.domAttrValue(n, "data-subtype")[1:]
		node.HeadingLevel, _ = strconv.Atoi(level)
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeBlockquote:
		content := strings.TrimSpace(lute.domText(n))
		if util.Caret == content {
			node.Type = ast.NodeText
			node.Tokens = []byte(content)
			tree.Context.Tip.AppendChild(node)
		}

		node.Type = ast.NodeBlockquote
		node.AppendChild(&ast.Node{Type: ast.NodeBlockquoteMarker, Tokens: []byte(">")})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeList:
		node.Type = ast.NodeList
		marker := lute.domAttrValue(n, "data-marker")
		node.ListData = &ast.ListData{}
		subType := lute.domAttrValue(n, "data-subtype")
		if "u" == subType {
			node.ListData.Typ = 0
		} else if "o" == subType {
			node.ListData.Typ = 1
		} else if "t" == subType {
			node.ListData.Typ = 3
		}
		node.ListData.Marker = []byte(marker)
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeListItem:
		marker := lute.domAttrValue(n, "data-marker")
		if ast.NodeList != tree.Context.Tip.Type {
			parent := &ast.Node{}
			parent.Type = ast.NodeList
			parent.ListData = &ast.ListData{}
			subType := lute.domAttrValue(n, "data-subtype")
			if "u" == subType {
				parent.ListData.Typ = 0
				parent.ListData.BulletChar = '*'
			} else if "o" == subType {
				parent.ListData.Typ = 1
				parent.ListData.Num, _ = strconv.Atoi(marker[:len(marker)-1])
				parent.ListData.Delimiter = '.'
			} else if "t" == subType {
				parent.ListData.Typ = 3
				parent.ListData.BulletChar = '*'
			}
			tree.Context.Tip.AppendChild(parent)
			tree.Context.Tip = parent
		}

		node.Type = ast.NodeListItem
		node.ListData = &ast.ListData{}
		subType := lute.domAttrValue(n, "data-subtype")
		if "u" == subType {
			node.ListData.Typ = 0
			node.ListData.BulletChar = '*'
		} else if "o" == subType {
			node.ListData.Typ = 1
			node.ListData.Num, _ = strconv.Atoi(marker[:len(marker)-1])
			node.ListData.Delimiter = '.'
		} else if "t" == subType {
			node.ListData.Typ = 3
			node.ListData.BulletChar = '*'
		}
		node.ListData.Marker = []byte(marker)
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeGitConflict:
		node.Type = ast.NodeGitConflict
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeSuperBlock:
		node.Type = ast.NodeSuperBlock
		tree.Context.Tip.AppendChild(node)
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockOpenMarker})
		layout := lute.domAttrValue(n, "data-sb-layout")
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockLayoutMarker, Tokens: []byte(layout)})
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeMathBlock:
		node.Type = ast.NodeMathBlock
		node.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker})
		content := lute.domAttrValue(n, "data-content")
		node.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent, Tokens: util.StrToBytes(content)})
		node.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker})
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeCodeBlock:
		node.Type = ast.NodeCodeBlock
		node.IsFencedCodeBlock = true
		node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: util.StrToBytes("```")})
		if language := lute.domAttrValue(n, "data-subtype"); "" != language {
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker, CodeBlockInfo: util.StrToBytes(language)})
			content := lute.domAttrValue(n, "data-content")
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockCode, Tokens: util.StrToBytes(content)})
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: util.StrToBytes("```")})
			tree.Context.Tip.AppendChild(node)
			return
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeHTMLBlock:
		node.Type = ast.NodeHTMLBlock
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeYamlFrontMatter:
		node.Type = ast.NodeYamlFrontMatter
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeThematicBreak:
		node.Type = ast.NodeThematicBreak
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeBlockEmbed:
		text := lute.domText(n)
		if "" == text {
			return
		}

		t := parse.Parse("", []byte(text), lute.ParseOptions)
		t.Root.LastChild.Unlink() // 移除 doc IAL
		if blockEmbed := t.Root.FirstChild; nil != blockEmbed && ast.NodeBlockEmbed == blockEmbed.Type {
			ial, id := node.KramdownIAL, node.ID
			node = blockEmbed
			node.KramdownIAL, node.ID = ial, id
			next := blockEmbed.Next
			tree.Context.Tip.AppendChild(node)
			appendNextToTip(next, tree)
			return
		}
		node.Type = ast.NodeText
		node.Tokens = []byte(text)
		tree.Context.Tip.AppendChild(node)
		defer tree.Context.ParentTip()
		return
	case ast.NodeIFrame:
		node.Type = ast.NodeIFrame
		n = lute.domChild(n.FirstChild, atom.Iframe)
		node.Tokens = lute.domHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeVideo:
		node.Type = ast.NodeVideo
		n = lute.domChild(n.FirstChild, atom.Video)
		node.Tokens = lute.domHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeAudio:
		node.Type = ast.NodeAudio
		n = lute.domChild(n.FirstChild, atom.Audio)
		node.Tokens = lute.domHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	default:
		switch n.DataAtom {
		case 0:
			node.Type = ast.NodeText
			node.Tokens = util.StrToBytes(n.Data)
			if ast.NodeDocument == tree.Context.Tip.Type {
				p := &ast.Node{Type: ast.NodeParagraph}
				tree.Context.Tip.AppendChild(p)
				tree.Context.Tip = p
			}
			lute.genASTContenteditable(n, tree)
			return
		case atom.U, atom.Code, atom.Strong, atom.Em, atom.Kbd, atom.Mark, atom.S, atom.Sub, atom.Sup, atom.Span:
			lute.genASTContenteditable(n, tree)
			return
		}

		if ast.NodeListItem == tree.Context.Tip.Type && atom.Input == n.DataAtom {
			node.Type = ast.NodeTaskListItemMarker
			node.TaskListItemChecked = lute.hasAttr(n, "checked")
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeInlineHTML
		node.Tokens = lute.domHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByBlockDOM(c, tree)
	}

	switch dataType {
	case ast.NodeSuperBlock:
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockCloseMarker})
	case ast.NodeCodeBlock:
		node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: util.StrToBytes("```")})
	}
}

func (lute *Lute) genASTContenteditable(n *html.Node, tree *parse.Tree) {
	if ast.NodeCodeBlock == tree.Context.Tip.Type {
		return
	}

	class := lute.domAttrValue(n, "class")
	if "svg" == class {
		return
	}

	content := n.Data
	node := &ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(content)}
	switch n.DataAtom {
	case 0:
		if "" == content {
			return
		}

		if ast.NodeLink == tree.Context.Tip.Type {
			node.Type = ast.NodeLinkText
		} else if ast.NodeHeading == tree.Context.Tip.Type {
			content = strings.ReplaceAll(content, "\n", "")
			node.Tokens = util.StrToBytes(content)
		}

		if lute.parentIs(n, atom.Table) {
			node.Tokens = util.StrToBytes(strings.ReplaceAll(content, "\n", "<br />"))
			array := lex.SplitWithoutBackslashEscape(node.Tokens, '|')
			node.Tokens = nil
			for i, tokens := range array {
				node.Tokens = append(node.Tokens, tokens...)
				if i < len(array)-1 {
					node.Tokens = append(node.Tokens, []byte("\\|")...)
				}
			}
		}
		if ast.NodeCodeSpan == tree.Context.Tip.Type || ast.NodeInlineMath == tree.Context.Tip.Type {
			tree.Context.Tip.FirstChild.Next.Tokens = util.StrToBytes(content)
			return
		}
		tree.Context.Tip.AppendChild(node)
	case atom.Thead:
		node.Type = ast.NodeTableHead
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Tbody:
	case atom.Tr:
		node.Type = ast.NodeTableRow
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Th, atom.Td:
		node.Type = ast.NodeTableCell
		align := lute.domAttrValue(n, "align")
		var tableAlign int
		switch align {
		case "left":
			tableAlign = 1
		case "center":
			tableAlign = 2
		case "right":
			tableAlign = 3
		default:
			tableAlign = 0
		}
		node.TableCellAlign = tableAlign
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Code:
		if lute.isEmptyText(n) {
			return
		}

		node.Type = ast.NodeCodeSpan
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanOpenMarker})
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanContent})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Span:
		dataType := lute.domAttrValue(n, "data-type")
		if "tag" == dataType {
			if nil == n.FirstChild {
				return
			}

			node.Type = ast.NodeTag
			node.AppendChild(&ast.Node{Type: ast.NodeTagOpenMarker})
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		} else if "inline-math" == dataType {
			node.Type = ast.NodeInlineMath
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathOpenMarker})
			content = lute.domAttrValue(n, "data-content")
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathContent, Tokens: util.StrToBytes(content)})
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathCloseMarker})
			tree.Context.Tip.AppendChild(node)
			return
		} else if "a" == dataType {
			if nil == n.FirstChild {
				// 丢弃没有锚文本的链接
				return
			}

			node.Type = ast.NodeLink
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		} else if "block-ref" == dataType {
			node.Type = ast.NodeBlockRef
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			id := lute.domAttrValue(n, "data-id")
			node.AppendChild(&ast.Node{Type: ast.NodeBlockRefID, Tokens: util.StrToBytes(id)})
			if refText := lute.domAttrValue(n, "data-anchor"); "" != refText {
				node.AppendChild(&ast.Node{Type: ast.NodeBlockRefSpace})
				refTextNode := &ast.Node{Type: ast.NodeBlockRefText, Tokens: util.StrToBytes(refText)}
				node.AppendChild(refTextNode)
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			tree.Context.Tip.AppendChild(node)
			return
		} else if "img" == dataType {
			img := n.FirstChild.FirstChild.NextSibling
			node.Type = ast.NodeImage
			node.AppendChild(&ast.Node{Type: ast.NodeBang})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			alt := lute.domAttrValue(img, "alt")
			node.AppendChild(&ast.Node{Type: ast.NodeLinkText, Tokens: util.StrToBytes(alt)})
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			src := lute.domAttrValue(img, "data-src")
			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: util.StrToBytes(src)})
			if title := lute.domAttrValue(img, "title"); "" != title {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
				node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: util.StrToBytes(title)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			tree.Context.Tip.AppendChild(node)
			lute.setSpanIAL(img, tree.Context.Tip.LastChild)
			return
		} else if "backslash" == dataType {
			node.Type = ast.NodeBackslash
			if nil == n.FirstChild {
				return
			}
			node.AppendChild(&ast.Node{Type: ast.NodeBackslashContent, Tokens: util.StrToBytes(n.FirstChild.NextSibling.Data)})
			tree.Context.Tip.AppendChild(node)
			return
		}
	case atom.Sub:
		if lute.isEmptyText(n) {
			return
		}

		node.Type = ast.NodeSub
		node.AppendChild(&ast.Node{Type: ast.NodeSubOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Sup:
		if lute.isEmptyText(n) {
			return
		}

		node.Type = ast.NodeSup
		node.AppendChild(&ast.Node{Type: ast.NodeSupOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.U:
		if lute.isEmptyText(n) {
			return
		}
		node.Type = ast.NodeUnderline
		node.AppendChild(&ast.Node{Type: ast.NodeUnderlineOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Kbd:
		if lute.isEmptyText(n) {
			return
		}
		node.Type = ast.NodeKbd
		node.AppendChild(&ast.Node{Type: ast.NodeKbdOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Br:
		if ast.NodeHeading == tree.Context.Tip.Type {
			return
		}

		node.Type = ast.NodeBr
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Em, atom.I:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if util.Caret == text {
			node.Tokens = util.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeEmphasis
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "*"
		}
		if "_" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeEmU8eOpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeEmA6kOpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && util.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.CaretTokens})
			if "_" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeEmU8eCloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, parse.Zwsp, "")

		// 开头结尾空格后会形成 * foo * 导致强调、加粗删除线标记失效，这里将空格移到右标记符前后 _*foo*_
		if strings.HasPrefix(n.FirstChild.Data, " ") && nil == n.FirstChild.PrevSibling {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}
		if strings.HasSuffix(n.FirstChild.Data, "\n") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, "\n")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: "\n"})
		}

		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Strong, atom.B:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if util.Caret == text {
			node.Tokens = util.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeStrong
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "**"
		}
		if "__" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongU8eOpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kOpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && util.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.CaretTokens})
			if "__" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeStrongU8eCloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, parse.Zwsp, "")
		if strings.HasPrefix(n.FirstChild.Data, " ") && nil == n.FirstChild.PrevSibling {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}
		if strings.HasSuffix(n.FirstChild.Data, "\n") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, "\n")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: "\n"})
		}

		lute.setSpanIAL(n, node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Del, atom.S, atom.Strike:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if util.Caret == text {
			node.Tokens = util.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeStrikethrough
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1OpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2OpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && util.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.CaretTokens})
			if "~" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1CloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, parse.Zwsp, "")
		if strings.HasPrefix(n.FirstChild.Data, " ") && nil == n.FirstChild.PrevSibling {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}
		if strings.HasSuffix(n.FirstChild.Data, "\n") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, "\n")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: "\n"})
		}

		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Mark:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if util.Caret == text {
			node.Tokens = util.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeMark
		marker := lute.domAttrValue(n, "data-marker")
		if "=" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeMark1OpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeMark2OpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && util.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.CaretTokens})
			if "=" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeMark1CloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, parse.Zwsp, "")
		if strings.HasPrefix(n.FirstChild.Data, " ") && nil == n.FirstChild.PrevSibling {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}
		if strings.HasSuffix(n.FirstChild.Data, "\n") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, "\n")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: "\n"})
		}

		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Img:
		if "emoji" == class {
			alt := lute.domAttrValue(n, "alt")
			node.Type = ast.NodeEmoji
			emojiImg := &ast.Node{Type: ast.NodeEmojiImg, Tokens: tree.EmojiImgTokens(alt, lute.domAttrValue(n, "src"))}
			emojiImg.AppendChild(&ast.Node{Type: ast.NodeEmojiAlias, Tokens: []byte(":" + alt + ":")})
			node.AppendChild(emojiImg)
			tree.Context.Tip.AppendChild(node)
			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTContenteditable(c, tree)
	}

	switch n.DataAtom {
	case atom.Code:
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanCloseMarker})
	case atom.Span:
		dataType := lute.domAttrValue(n, "data-type")
		if "tag" == dataType {
			node.AppendChild(&ast.Node{Type: ast.NodeTagCloseMarker})
		} else if "a" == dataType {
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			href := lute.domAttrValue(n, "data-href")
			if "" != lute.RenderOptions.LinkBase {
				href = strings.ReplaceAll(href, lute.RenderOptions.LinkBase, "")
			}
			if "" != lute.RenderOptions.LinkPrefix {
				href = strings.ReplaceAll(href, lute.RenderOptions.LinkPrefix, "")
			}
			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: []byte(href)})
			linkTitle := lute.domAttrValue(n, "data-title")
			if "" != linkTitle {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
				node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: []byte(linkTitle)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
		} else if "block-ref" == dataType {
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
		}
	case atom.Sub:
		node.AppendChild(&ast.Node{Type: ast.NodeSubCloseMarker})
	case atom.Sup:
		node.AppendChild(&ast.Node{Type: ast.NodeSupCloseMarker})
	case atom.U:
		node.AppendChild(&ast.Node{Type: ast.NodeUnderlineCloseMarker})
	case atom.Kbd:
		node.AppendChild(&ast.Node{Type: ast.NodeKbdCloseMarker})
	case atom.Em, atom.I:
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "*"
		}
		if "_" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeEmU8eCloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: []byte(marker)})
		}
	case atom.Strong, atom.B:
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "**"
		}
		if "__" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongU8eCloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker, Tokens: []byte(marker)})
		}
	case atom.Del, atom.S, atom.Strike:
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1CloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker, Tokens: []byte(marker)})
		}
	case atom.Mark:
		marker := lute.domAttrValue(n, "data-marker")
		if "=" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeMark1CloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker, Tokens: []byte(marker)})
		}
	}
}

func (lute *Lute) setSpanIAL(n *html.Node, node *ast.Node) {
	insertedIAL := false
	if style := lute.domAttrValue(n, "style"); "" != style {
		node.SetIALAttr("style", style)
		ialTokens := parse.IAL2Tokens(node.KramdownIAL)
		ial := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
		node.InsertAfter(ial)
		insertedIAL = true
	}

	if nil != n.Parent && nil != n.Parent.Parent {
		if parentStyle := lute.domAttrValue(n.Parent.Parent, "style"); "" != parentStyle {
			if insertedIAL {
				m := parse.Tokens2IAL(node.Next.Tokens)
				m = append(m, []string{"parent-style", parentStyle})
				node.Next.Tokens = parse.IAL2Tokens(m)
				node.SetIALAttr("parent-style", parentStyle)
				node.KramdownIAL = m
			} else {
				node.SetIALAttr("parent-style", parentStyle)
				ialTokens := parse.IAL2Tokens(node.KramdownIAL)
				ial := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
				node.InsertAfter(ial)
			}
		}
	}
}

func (lute *Lute) setBlockIAL(n *html.Node, node *ast.Node) (ialTokens []byte) {
	ialTokens = []byte("{: id=\"" + node.ID + "\"")

	if refcount := lute.domAttrValue(n, "refcount"); "" != refcount {
		refcount = html.UnescapeHTMLStr(refcount)
		refcount = html.EscapeHTMLStr(refcount)
		node.SetIALAttr("refcount", refcount)
		ialTokens = append(ialTokens, []byte(" refcount=\""+refcount+"\"")...)
	}

	if bookmark := lute.domAttrValue(n, "bookmark"); "" != bookmark {
		bookmark = html.UnescapeHTMLStr(bookmark)
		bookmark = html.EscapeHTMLStr(bookmark)
		node.SetIALAttr("bookmark", bookmark)
		ialTokens = append(ialTokens, []byte(" bookmark=\""+bookmark+"\"")...)
	}

	if style := lute.domAttrValue(n, "style"); "" != style {
		node.SetIALAttr("style", style)
		ialTokens = append(ialTokens, []byte(" style=\""+style+"\"")...)
	}

	if name := lute.domAttrValue(n, "name"); "" != name {
		name = html.UnescapeHTMLStr(name)
		name = html.EscapeHTMLStr(name)
		node.SetIALAttr("name", name)
		ialTokens = append(ialTokens, []byte(" name=\""+name+"\"")...)
	}

	if memo := lute.domAttrValue(n, "memo"); "" != memo {
		memo = html.UnescapeHTMLStr(memo)
		memo = html.EscapeHTMLStr(memo)
		node.SetIALAttr("memo", memo)
		ialTokens = append(ialTokens, []byte(" memo=\""+memo+"\"")...)
	}

	if alias := lute.domAttrValue(n, "alias"); "" != alias {
		alias = html.UnescapeHTMLStr(alias)
		alias = html.EscapeHTMLStr(alias)
		node.SetIALAttr("alias", alias)
		ialTokens = append(ialTokens, []byte(" alias=\""+alias+"\"")...)
	}

	if fold := lute.domAttrValue(n, "fold"); "" != fold {
		fold = html.UnescapeHTMLStr(fold)
		fold = html.EscapeHTMLStr(fold)
		node.SetIALAttr("fold", fold)
		ialTokens = append(ialTokens, []byte(" fold=\""+fold+"\"")...)
	}

	if headingFold := lute.domAttrValue(n, "heading-fold"); "" != headingFold {
		headingFold = html.UnescapeHTMLStr(headingFold)
		headingFold = html.EscapeHTMLStr(headingFold)
		node.SetIALAttr("heading-fold", headingFold)
		ialTokens = append(ialTokens, []byte(" heading-fold=\""+headingFold+"\"")...)
	}

	if parentFold := lute.domAttrValue(n, "parent-fold"); "" != parentFold {
		parentFold = html.UnescapeHTMLStr(parentFold)
		parentFold = html.EscapeHTMLStr(parentFold)
		node.SetIALAttr("parent-fold", parentFold)
		ialTokens = append(ialTokens, []byte(" parent-fold=\""+parentFold+"\"")...)
	}

	if updated := lute.domAttrValue(n, "updated"); "" != updated {
		updated = html.UnescapeHTMLStr(updated)
		updated = html.EscapeHTMLStr(updated)
		node.SetIALAttr("updated", updated)
		ialTokens = append(ialTokens, []byte(" updated=\""+updated+"\"")...)
	}

	if linewrap := lute.domAttrValue(n, "linewrap"); "" != linewrap {
		linewrap = html.UnescapeHTMLStr(linewrap)
		linewrap = html.EscapeHTMLStr(linewrap)
		node.SetIALAttr("linewrap", linewrap)
		ialTokens = append(ialTokens, []byte(" linewrap=\""+linewrap+"\"")...)
	}

	if ligatures := lute.domAttrValue(n, "ligatures"); "" != ligatures {
		ligatures = html.UnescapeHTMLStr(ligatures)
		ligatures = html.EscapeHTMLStr(ligatures)
		node.SetIALAttr("ligatures", ligatures)
		ialTokens = append(ialTokens, []byte(" ligatures=\""+ligatures+"\"")...)
	}

	if linenumber := lute.domAttrValue(n, "linenumber"); "" != linenumber {
		linenumber = html.UnescapeHTMLStr(linenumber)
		linenumber = html.EscapeHTMLStr(linenumber)
		node.SetIALAttr("linenumber", linenumber)
		ialTokens = append(ialTokens, []byte(" linenumber=\""+linenumber+"\"")...)
	}

	if customAttrs := lute.domCustomAttrs(n); nil != customAttrs {
		for k, v := range customAttrs {
			v = html.UnescapeHTMLStr(v)
			v = html.EscapeHTMLStr(v)
			node.SetIALAttr(k, v)
			ialTokens = append(ialTokens, []byte(" "+k+"=\""+v+"\"")...)
		}
	}

	ialTokens = append(ialTokens, '}')
	return ialTokens
}

func appendNextToTip(next *ast.Node, tree *parse.Tree) {
	var nodes []*ast.Node
	for n := next; nil != n; n = n.Next {
		nodes = append(nodes, n)
	}
	for _, n := range nodes {
		tree.Context.Tip.AppendChild(n)
	}
}
