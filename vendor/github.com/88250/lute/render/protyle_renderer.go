// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// BlockRenderer 描述了 WYSIWYG Block DOM 渲染器。
type BlockRenderer struct {
	*BaseRenderer
	NodeIndex int
}

// NewBlockRenderer 创建一个 WYSIWYG Block DOM 渲染器。
func NewBlockRenderer(tree *parse.Tree, options *Options) *BlockRenderer {
	ret := &BlockRenderer{BaseRenderer: NewBaseRenderer(tree, options), NodeIndex: options.NodeIndexStart}
	ret.RendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.RendererFuncs[ast.NodeText] = ret.renderText
	ret.RendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	ret.RendererFuncs[ast.NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.RendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.RendererFuncs[ast.NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.RendererFuncs[ast.NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	ret.RendererFuncs[ast.NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	ret.RendererFuncs[ast.NodeCodeBlockCode] = ret.renderCodeBlockCode
	ret.RendererFuncs[ast.NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.RendererFuncs[ast.NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	ret.RendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent
	ret.RendererFuncs[ast.NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.RendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	ret.RendererFuncs[ast.NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	ret.RendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent
	ret.RendererFuncs[ast.NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	ret.RendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	ret.RendererFuncs[ast.NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	ret.RendererFuncs[ast.NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	ret.RendererFuncs[ast.NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	ret.RendererFuncs[ast.NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.RendererFuncs[ast.NodeStrong] = ret.renderStrong
	ret.RendererFuncs[ast.NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	ret.RendererFuncs[ast.NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	ret.RendererFuncs[ast.NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	ret.RendererFuncs[ast.NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.RendererFuncs[ast.NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.RendererFuncs[ast.NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	ret.RendererFuncs[ast.NodeHeadingID] = ret.renderHeadingID
	ret.RendererFuncs[ast.NodeList] = ret.renderList
	ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.RendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.RendererFuncs[ast.NodeLink] = ret.renderLink
	ret.RendererFuncs[ast.NodeImage] = ret.renderImage
	ret.RendererFuncs[ast.NodeBang] = ret.renderBang
	ret.RendererFuncs[ast.NodeOpenBracket] = ret.renderOpenBracket
	ret.RendererFuncs[ast.NodeCloseBracket] = ret.renderCloseBracket
	ret.RendererFuncs[ast.NodeOpenParen] = ret.renderOpenParen
	ret.RendererFuncs[ast.NodeCloseParen] = ret.renderCloseParen
	ret.RendererFuncs[ast.NodeOpenBrace] = ret.renderOpenBrace
	ret.RendererFuncs[ast.NodeCloseBrace] = ret.renderCloseBrace
	ret.RendererFuncs[ast.NodeLinkText] = ret.renderLinkText
	ret.RendererFuncs[ast.NodeLinkSpace] = ret.renderLinkSpace
	ret.RendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	ret.RendererFuncs[ast.NodeLinkTitle] = ret.renderLinkTitle
	ret.RendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	ret.RendererFuncs[ast.NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	ret.RendererFuncs[ast.NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	ret.RendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.RendererFuncs[ast.NodeTable] = ret.renderTable
	ret.RendererFuncs[ast.NodeTableHead] = ret.renderTableHead
	ret.RendererFuncs[ast.NodeTableRow] = ret.renderTableRow
	ret.RendererFuncs[ast.NodeTableCell] = ret.renderTableCell
	ret.RendererFuncs[ast.NodeEmoji] = ret.renderEmoji
	ret.RendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.RendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg
	ret.RendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias
	ret.RendererFuncs[ast.NodeFootnotesDefBlock] = ret.renderFootnotesDefBlock
	ret.RendererFuncs[ast.NodeFootnotesDef] = ret.renderFootnotesDef
	ret.RendererFuncs[ast.NodeFootnotesRef] = ret.renderFootnotesRef
	ret.RendererFuncs[ast.NodeToC] = ret.renderToC
	ret.RendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	ret.RendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	ret.RendererFuncs[ast.NodeHTMLEntity] = ret.renderHtmlEntity
	ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter
	ret.RendererFuncs[ast.NodeYamlFrontMatterOpenMarker] = ret.renderYamlFrontMatterOpenMarker
	ret.RendererFuncs[ast.NodeYamlFrontMatterContent] = ret.renderYamlFrontMatterContent
	ret.RendererFuncs[ast.NodeYamlFrontMatterCloseMarker] = ret.renderYamlFrontMatterCloseMarker
	ret.RendererFuncs[ast.NodeBlockRef] = ret.renderBlockRef
	ret.RendererFuncs[ast.NodeBlockRefID] = ret.renderBlockRefID
	ret.RendererFuncs[ast.NodeBlockRefSpace] = ret.renderBlockRefSpace
	ret.RendererFuncs[ast.NodeBlockRefText] = ret.renderBlockRefText
	ret.RendererFuncs[ast.NodeBlockEmbed] = ret.renderNodeBlockEmbed
	ret.RendererFuncs[ast.NodeBlockRefTextTplRenderResult] = ret.renderBlockRefTextTplRenderResult
	ret.RendererFuncs[ast.NodeMark] = ret.renderMark
	ret.RendererFuncs[ast.NodeMark1OpenMarker] = ret.renderMark1OpenMarker
	ret.RendererFuncs[ast.NodeMark1CloseMarker] = ret.renderMark1CloseMarker
	ret.RendererFuncs[ast.NodeMark2OpenMarker] = ret.renderMark2OpenMarker
	ret.RendererFuncs[ast.NodeMark2CloseMarker] = ret.renderMark2CloseMarker
	ret.RendererFuncs[ast.NodeSup] = ret.renderSup
	ret.RendererFuncs[ast.NodeSupOpenMarker] = ret.renderSupOpenMarker
	ret.RendererFuncs[ast.NodeSupCloseMarker] = ret.renderSupCloseMarker
	ret.RendererFuncs[ast.NodeSub] = ret.renderSub
	ret.RendererFuncs[ast.NodeSubOpenMarker] = ret.renderSubOpenMarker
	ret.RendererFuncs[ast.NodeSubCloseMarker] = ret.renderSubCloseMarker
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.RendererFuncs[ast.NodeKramdownSpanIAL] = ret.renderKramdownSpanIAL
	ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	ret.RendererFuncs[ast.NodeBlockQueryEmbedScript] = ret.renderBlockQueryEmbedScript
	ret.RendererFuncs[ast.NodeTag] = ret.renderTag
	ret.RendererFuncs[ast.NodeTagOpenMarker] = ret.renderTagOpenMarker
	ret.RendererFuncs[ast.NodeTagCloseMarker] = ret.renderTagCloseMarker
	ret.RendererFuncs[ast.NodeLinkRefDefBlock] = ret.renderLinkRefDefBlock
	ret.RendererFuncs[ast.NodeLinkRefDef] = ret.renderLinkRefDef
	ret.RendererFuncs[ast.NodeSuperBlock] = ret.renderSuperBlock
	ret.RendererFuncs[ast.NodeSuperBlockOpenMarker] = ret.renderSuperBlockOpenMarker
	ret.RendererFuncs[ast.NodeSuperBlockLayoutMarker] = ret.renderSuperBlockLayoutMarker
	ret.RendererFuncs[ast.NodeSuperBlockCloseMarker] = ret.renderSuperBlockCloseMarker
	ret.RendererFuncs[ast.NodeGitConflict] = ret.renderGitConflict
	ret.RendererFuncs[ast.NodeGitConflictOpenMarker] = ret.renderGitConflictOpenMarker
	ret.RendererFuncs[ast.NodeGitConflictContent] = ret.renderGitConflictContent
	ret.RendererFuncs[ast.NodeGitConflictCloseMarker] = ret.renderGitConflictCloseMarker
	ret.RendererFuncs[ast.NodeIFrame] = ret.renderIFrame
	ret.RendererFuncs[ast.NodeVideo] = ret.renderVideo
	ret.RendererFuncs[ast.NodeAudio] = ret.renderAudio
	ret.RendererFuncs[ast.NodeKbd] = ret.renderKbd
	ret.RendererFuncs[ast.NodeKbdOpenMarker] = ret.renderKbdOpenMarker
	ret.RendererFuncs[ast.NodeKbdCloseMarker] = ret.renderKbdCloseMarker
	ret.RendererFuncs[ast.NodeUnderline] = ret.renderUnderline
	ret.RendererFuncs[ast.NodeUnderlineOpenMarker] = ret.renderUnderlineOpenMarker
	ret.RendererFuncs[ast.NodeUnderlineCloseMarker] = ret.renderUnderlineCloseMarker
	ret.RendererFuncs[ast.NodeBr] = ret.renderBr
	return ret
}

func (r *BlockRenderer) renderBr(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<br />")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderUnderline(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderUnderlineOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<u>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderUnderlineCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</u>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderKbd(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderKbdOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<kbd>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderKbdCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</kbd>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		tokens := node.ChildByType(ast.NodeBlockQueryEmbedScript).Tokens
		tokens = html.EscapeHTML(bytes.ReplaceAll(tokens, util.CaretTokens, nil))
		attrs = append(attrs, []string{"data-content", util.BytesToStr(tokens)})
		r.blockNodeAttrs(node, &attrs, "render-node")
		r.Tag("div", attrs, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockQueryEmbedScript(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderVideo(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "iframe")
		r.Tag("div", attrs, false)
		r.Tag("div", [][]string{{"class", "iframe-content"}}, false)
		tokens := bytes.ReplaceAll(node.Tokens, util.CaretTokens, nil)
		dataSrc := r.tagSrc(tokens)
		src := r.LinkPath(dataSrc)
		tokens = r.replaceSrc(tokens, src, dataSrc)
		r.Write(tokens)
	} else {
		r.Tag("span", [][]string{{"class", "protyle-action__drag"}, {"contenteditable", "false"}}, false)
		r.Tag("/span", nil, false)
		r.WriteString(parse.Zwsp)
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderAudio(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "iframe")
		r.Tag("div", attrs, false)
		r.Tag("div", [][]string{{"class", "iframe-content"}}, false)
		tokens := bytes.ReplaceAll(node.Tokens, util.CaretTokens, nil)
		dataSrc := r.tagSrc(tokens)
		src := r.LinkPath(dataSrc)
		tokens = r.replaceSrc(tokens, src, dataSrc)
		r.Write(tokens)
		r.WriteString(parse.Zwsp)
	} else {
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderIFrame(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "iframe")
		r.Tag("div", attrs, false)
		r.Tag("div", [][]string{{"class", "iframe-content"}}, false)
		tokens := bytes.ReplaceAll(node.Tokens, util.CaretTokens, nil)
		dataSrc := r.tagSrc(tokens)
		src := r.LinkPath(dataSrc)
		tokens = r.replaceSrc(tokens, src, dataSrc)
		r.Write(tokens)
	} else {
		r.Tag("span", [][]string{{"class", "protyle-action__drag"}, {"contenteditable", "false"}}, false)
		r.Tag("/span", nil, false)
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) replaceSrc(tokens, src, dataSrc []byte) []byte {
	src1 := append([]byte(" src=\""), src...)
	src1 = append(src1, []byte("\"")...)
	dataSrc1 := append([]byte(" src=\""), dataSrc...)
	dataSrc1 = append(dataSrc1, []byte("\"")...)
	if bytes.Contains(tokens, []byte("data-src=")) {
		return bytes.ReplaceAll(tokens, dataSrc1, src1)
	}
	src1 = append(src1, []byte(" data-src=\""+util.BytesToStr(dataSrc)+"\"")...)
	return bytes.ReplaceAll(tokens, dataSrc1, src1)
}

func (r *BlockRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		idNode := node.ChildByType(ast.NodeBlockRefID)
		var anchor string
		if refTextNode := node.ChildByType(ast.NodeBlockRefText); nil != refTextNode {
			anchor = strings.ReplaceAll(refTextNode.Text(), util.Caret, "")
		}
		attrs := [][]string{{"data-type", "block-ref"}, {"data-id", idNode.TokensStr()}, {"data-anchor", anchor}}
		r.Tag("span", attrs, false)
		refTextNode := node.ChildByType(ast.NodeBlockRefTextTplRenderResult)
		var refText string
		if nil != refTextNode {
			refText = refTextNode.TokensStr()
		} else {
			refText = anchor
		}
		r.WriteString(refText)
		r.Tag("/span", nil, false)
		return ast.WalkSkipChildren
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockRefText(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockRefTextTplRenderResult(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderGitConflictCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderGitConflictContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "git-conflict")
		r.Tag("div", attrs, false)
		attrs = [][]string{{"contenteditable", "false"}, {"spellcheck", "false"}}
		r.Tag("div", attrs, false)

		tokens := bytes.TrimSpace(node.Tokens)
		r.Write(html.EscapeHTML(tokens))
	} else {
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}

	return ast.WalkContinue
}

func (r *BlockRenderer) renderGitConflictOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderGitConflict(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTag(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTagOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		content := node.Parent.Text()
		content = strings.ReplaceAll(content, util.Caret, "")
		r.Tag("span", [][]string{{"data-type", "tag"}, {"data-content", content}}, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSuperBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "sb")
		layout := node.FirstChild.Next.TokensStr()
		if "" == layout {
			layout = "row"
		}
		attrs = append(attrs, []string{"data-sb-layout", layout})
		r.Tag("div", attrs, false)
	} else {
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSuperBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSuperBlockLayoutMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSuperBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div data-block=\"0\" data-type=\"link-ref-defs-block\">")
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		dest := node.FirstChild.ChildByType(ast.NodeLinkDest).Tokens
		destStr := util.BytesToStr(dest)
		r.WriteString("[" + util.BytesToStr(node.Tokens) + "]:")
		if util.Caret != destStr {
			r.WriteString(" ")
		}
		r.WriteString(destStr + "\n")
	}
	return ast.WalkSkipChildren
}

func (r *BlockRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previewTokens := bytes.TrimSpace(node.Tokens)
		codeLen := len(previewTokens)
		codeIsEmpty := 1 > codeLen || (len(util.Caret) == codeLen && util.Caret == string(node.Tokens))
		r.Tag("pre", nil, false)
		r.Tag("code", [][]string{{"data-type", "yaml-front-matter"}}, false)
		if codeIsEmpty {
			r.WriteString(util.FrontEndCaret + "\n")
		} else {
			r.Write(html.EscapeHTML(previewTokens))
		}
		r.WriteString("</code></pre>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString(`<div class="protyle-wysiwyg__block" data-type="yaml-front-matter" data-block="0">`)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	r.WriteString("<span class=\"protyle-wysiwyg__block\" data-type=\"html-entity\">")
	r.Tag("code", [][]string{{"data-type", "html-entity"}, {"style", "display: none"}}, false)
	tokens := node.HtmlEntityTokens
	r.Write(html.EscapeHTML(tokens))
	r.WriteString("</code>")

	r.Tag("span", [][]string{{"class", "protyle-wysiwyg__preview"}, {"data-render", "2"}}, false)
	r.Tag("code", nil, false)
	previewTokens := bytes.ReplaceAll(node.HtmlEntityTokens, util.CaretTokens, nil)
	r.Write(previewTokens)
	r.Tag("/code", nil, false)
	r.Tag("/span", nil, false)
	r.WriteString("</span>")
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<span data-type=\"backslash\">")
		r.WriteString("<span>")
		r.WriteByte(lex.ItemBackslash)
		r.WriteString("</span>")
	} else {
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	return r.BaseRenderer.renderToC(node, entering)
}

func (r *BlockRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div data-block=\"0\" data-type=\"footnotes-block\">")
		r.WriteString("<ol data-type=\"footnotes-defs-ol\">")
	} else {
		r.WriteString("</ol></div>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if r.RenderingFootnotes {
		return ast.WalkContinue
	}

	if entering {
		if nil != node.Previous && bytes.EqualFold(node.Previous.Tokens, node.Tokens) {
			return ast.WalkContinue
		}

		r.WriteString("<li data-type=\"footnotes-li\" data-marker=\"" + string(node.Tokens) + "\">")
		for c := node.FirstChild; nil != c; c = c.Next {
			ast.Walk(c, func(n *ast.Node, entering bool) ast.WalkStatus {
				return r.RendererFuncs[n.Type](n, entering)
			})
		}
		r.WriteString("</li>")
		return ast.WalkSkipChildren
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		idx, def := r.Tree.FindFootnotesDef(node.Tokens)
		idxStr := strconv.Itoa(idx)
		label := def.Text()
		r.Tag("sup", [][]string{{"data-type", "footnotes-ref"}, {"data-footnotes-label", string(node.FootnotesRefLabel)},
			{"class", "protyle-tooltipped protyle-tooltipped__s"}, {"aria-label", SubStr(html.EscapeHTMLStr(label), 24)}}, false)
		r.WriteString(idxStr)
		r.WriteString("</sup>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	noHighlight := false
	var language string
	if nil != node.FirstChild.Next && 0 < len(node.FirstChild.Next.CodeBlockInfo) {
		language = util.BytesToStr(node.FirstChild.Next.CodeBlockInfo)
		noHighlight = r.NoHighlight(language)
	}

	if entering {
		if noHighlight {
			var attrs [][]string
			r.blockNodeAttrs(node, &attrs, "render-node")
			tokens := html.EscapeHTML(node.FirstChild.Next.Next.Tokens)
			tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
			tokens = bytes.TrimSpace(tokens)
			attrs = append(attrs, []string{"data-content", util.BytesToStr(tokens)})
			attrs = append(attrs, []string{"data-subtype", language})
			r.Tag("div", attrs, false)
			r.Tag("div", [][]string{{"spin", "1"}}, false)
			r.Tag("/div", nil, false)
			r.renderIAL(node)
			r.Tag("/div", nil, false)
			return ast.WalkSkipChildren
		}

		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "code-block")
		r.Tag("div", attrs, false)
	} else {
		if noHighlight {
			return ast.WalkSkipChildren
		}
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	r.Tag("div", [][]string{{"class", "protyle-action protyle-icons"}}, false)
	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(util.Caret) == codeLen && util.Caret == string(node.Tokens))
	var language string
	caretInInfo := false
	if nil != node.Previous {
		caretInInfo = bytes.Contains(node.Previous.CodeBlockInfo, util.CaretTokens)
		node.Previous.CodeBlockInfo = bytes.ReplaceAll(node.Previous.CodeBlockInfo, util.CaretTokens, nil)
	}

	attrs := [][]string{{"class", "protyle-action__language"}, {"contenteditable", "false"}}
	if nil != node.Previous && 0 < len(node.Previous.CodeBlockInfo) {
		infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
		language = string(infoWords[0])
	}

	r.Tag("span", attrs, false)
	r.WriteString(language)
	r.Tag("/span", nil, false)

	r.Tag("span", [][]string{{"class", "protyle-action__copy"}}, false)
	r.Tag("/span", nil, false)
	r.Tag("/div", nil, false)

	attrs = [][]string{}
	r.contenteditable(&attrs)
	r.spellcheck(&attrs)
	r.Tag("div", attrs, false)
	if codeIsEmpty {
		if caretInInfo {
			r.WriteString(util.FrontEndCaret)
		}
	} else {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	r.Tag("/div", nil, false)
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.Options.AutoSpace {
			if text := node.PreviousNodeText(); "" != text {
				lastc, _ := utf8.DecodeLastRuneInString(text)
				if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
					r.WriteByte(lex.ItemSpace)
				}
			}
		}
	} else {
		if r.Options.AutoSpace {
			if text := node.NextNodeText(); "" != text {
				firstc, _ := utf8.DecodeRuneInString(text)
				if unicode.IsLetter(firstc) || unicode.IsDigit(firstc) {
					r.WriteByte(lex.ItemSpace)
				}
			}
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.EscapeHTML(node.Next.Tokens)
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
		r.Tag("span", [][]string{{"data-type", "inline-math"}, {"data-subtype", "math"}, {"data-content", util.BytesToStr(tokens)}, {"contenteditable", "false"}, {"class", "render-node"}}, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/span", nil, false)
		if bytes.Contains(node.Previous.Tokens, util.CaretTokens) {
			r.WriteString(util.Caret)
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	var attrs [][]string
	r.blockNodeAttrs(node, &attrs, "render-node")
	tokens := html.EscapeHTML(node.FirstChild.Next.Tokens)
	tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
	tokens = bytes.TrimSpace(tokens)
	attrs = append(attrs, []string{"data-content", util.BytesToStr(tokens)})
	attrs = append(attrs, []string{"data-subtype", "math"})
	r.Tag("div", attrs, false)
	r.Tag("div", [][]string{{"spin", "1"}}, false)
	r.Tag("/div", nil, false)
	r.renderIAL(node)
	r.Tag("/div", nil, false)
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	tag := "td"
	if ast.NodeTableHead == node.Parent.Parent.Type {
		tag = "th"
	}
	if entering {
		var attrs [][]string
		switch node.TableCellAlign {
		case 1:
			attrs = append(attrs, []string{"align", "left"})
		case 2:
			attrs = append(attrs, []string{"align", "center"})
		case 3:
			attrs = append(attrs, []string{"align", "right"})
		}
		r.Tag(tag, attrs, false)
	} else {
		r.Tag("/"+tag, nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("tr", nil, false)
	} else {
		r.Tag("/tr", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("thead", nil, false)
	} else {
		r.Tag("/thead", nil, false)
		if nil != node.Next {
			r.Tag("tbody", nil, false)
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if nil == node.FirstChild {
		return ast.WalkSkipChildren
	}

	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "table")
		r.Tag("div", attrs, false)
		attrs = [][]string{}
		r.contenteditable(&attrs)
		r.spellcheck(&attrs)
		r.Tag("div", attrs, false)
		r.Tag("table", nil, false)
	} else {
		if nil != node.FirstChild.Next {
			r.Tag("/tbody", nil, false)
		}
		r.Tag("/table", nil, false)
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("s", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("s", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeImage != node.Parent.Type {
			r.Write(node.Tokens)
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		attrs := [][]string{{"contenteditable", "false"}, {"data-type", "img"}, {"class", "img"}}
		parentStyle := node.IALAttr("parent-style")
		if "" != parentStyle { // 手动设置了位置
			attrs = append(attrs, []string{"style", parentStyle})
		}
		r.Tag("span", attrs, false)
		r.Tag("span", nil, false)
		r.Tag("span", [][]string{{"class", "protyle-action protyle-icons"}}, false)
		r.WriteString("<span><svg class=\"svg\"><use xlink:href=\"#iconMore\"></use></svg></span>")
		r.Tag("/span", nil, false)
	} else {
		destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
		destTokens = bytes.ReplaceAll(destTokens, util.CaretTokens, nil)
		dataSrc := util.BytesToStr(destTokens)
		src := util.BytesToStr(r.LinkPath(destTokens))
		attrs := [][]string{{"src", src}, {"data-src", dataSrc}}
		alt := node.ChildByType(ast.NodeLinkText)
		if nil != alt && 0 < len(alt.Tokens) {
			attrs = append(attrs, []string{"alt", util.BytesToStr(alt.Tokens)})
		}

		title := node.ChildByType(ast.NodeLinkTitle)
		var titleTokens []byte
		if nil != title && 0 < len(title.Tokens) {
			titleTokens = title.Tokens
			attrs = append(attrs, []string{"title", util.BytesToStr(titleTokens)})
		}

		attrs = append(attrs, r.NodeAttrs(node.Parent)...)
		if style := node.IALAttr("style"); "" != style {
			attrs = append(attrs, []string{"style", style})
		}
		r.Tag("img", attrs, true)
		// XSS 过滤
		buf := r.Writer.Bytes()
		idx := bytes.LastIndex(buf, []byte("<img src="))
		imgBuf := buf[idx:]
		imgBuf = r.tagSrcPath(imgBuf)
		r.Writer.Truncate(idx)
		r.Writer.Write(imgBuf)

		r.Tag("span", [][]string{{"class", "protyle-action__drag"}}, false)
		r.Tag("/span", nil, false)

		attrs = [][]string{{"class", "protyle-action__title"}}
		r.Tag("span", attrs, false)
		r.Writer.Write(titleTokens)
		r.Tag("/span", nil, false)
		r.Tag("/span", nil, false)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		dest := node.ChildByType(ast.NodeLinkDest)
		destTokens := dest.Tokens
		destTokens = r.LinkPath(destTokens)
		caretInDest := bytes.Contains(destTokens, util.CaretTokens)
		if caretInDest {
			text := node.ChildByType(ast.NodeLinkText)
			text.Tokens = append(text.Tokens, util.CaretTokens...)
			destTokens = bytes.ReplaceAll(destTokens, util.CaretTokens, nil)
		}
		attrs := [][]string{{"data-type", "a"}, {"data-href", string(destTokens)}}
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			attrs = append(attrs, []string{"data-title", util.BytesToStr(title.Tokens)})
		}
		r.Tag("span", attrs, false)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	// Protyle 中不存在 HTML 块，所以应该不会进入这里的渲染实现

	if entering {
		var attrs [][]string
		node.Type = ast.NodeParagraph
		r.blockNodeAttrs(node, &attrs, "p")
		r.Tag("div", attrs, false)
		attrs = [][]string{{"contenteditable", "false"}, {"spellcheck", "false"}}
		r.Tag("div", attrs, false)

		tokens := bytes.TrimSpace(node.Tokens)
		r.Write(html.EscapeHTML(tokens))
	} else {
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	if bytes.Equal(node.Tokens, []byte("<br />")) && node.ParentIs(ast.NodeTableCell) {
		r.Write(node.Tokens)
		return ast.WalkContinue
	}

	if bytes.Equal(node.Tokens, []byte("<u>")) || bytes.Equal(node.Tokens, []byte("</u>")) {
		r.Write(node.Tokens)
		return ast.WalkContinue
	}

	r.Tag("code", [][]string{{"data-type", "html-inline"}}, false)
	tokens := html.EscapeHTML(node.Tokens)
	r.Write(tokens)
	r.WriteString("</code>")
	return ast.WalkContinue
}

func (r *BlockRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderNodeBlockEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	// 嵌入节点 !((id)) 已被废弃，这里按照段落渲染
	if entering {
		var attrs [][]string
		node.Type = ast.NodeParagraph
		r.blockNodeAttrs(node, &attrs, "p")
		r.Tag("div", attrs, false)
		attrs = [][]string{}
		r.contenteditable(&attrs)
		r.spellcheck(&attrs)
		r.Tag("div", attrs, false)
		idNode := node.ChildByType(ast.NodeBlockEmbedID)
		id := idNode.TokensStr()
		r.WriteString("!((" + id + "))")
	} else {
		if (nil != node.LastChild && util.Caret == node.LastChild.TokensStr() && nil != node.LastChild.Previous && ast.NodeImage == node.LastChild.Previous.Type) ||
			(nil != node.LastChild && ast.NodeImage == node.LastChild.Type) ||
			(nil != node.LastChild && ast.NodeKramdownSpanIAL == node.LastChild.Type && nil != node.LastChild.Previous && ast.NodeImage == node.LastChild.Previous.Type) {
			r.WriteString("\n") // 主要是为了解决 img 插入符后置问题
		}
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkSkipChildren
}

func (r *BlockRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "p")
		r.Tag("div", attrs, false)
		attrs = [][]string{}
		r.contenteditable(&attrs)
		r.spellcheck(&attrs)
		r.Tag("div", attrs, false)
	} else {
		if (nil != node.LastChild && util.Caret == node.LastChild.TokensStr() && nil != node.LastChild.Previous && ast.NodeImage == node.LastChild.Previous.Type) ||
			(nil != node.LastChild && ast.NodeImage == node.LastChild.Type) ||
			(nil != node.LastChild && ast.NodeKramdownSpanIAL == node.LastChild.Type && nil != node.LastChild.Previous && ast.NodeImage == node.LastChild.Previous.Type) {
			r.WriteString("\n") // 主要是为了解决 img 插入符后置问题
		}
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var tokens []byte
		if r.Options.AutoSpace && ast.NodeKbd != node.Parent.Type {
			tokens = r.Space(node.Tokens)
		} else {
			tokens = node.Tokens
		}
		r.Write(html.EscapeHTML(tokens))
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.Options.AutoSpace {
			if text := node.PreviousNodeText(); "" != text {
				lastc, _ := utf8.DecodeLastRuneInString(text)
				if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
					r.WriteByte(lex.ItemSpace)
				}
			}
		}
	} else {
		if r.Options.AutoSpace {
			if text := node.NextNodeText(); "" != text {
				firstc, _ := utf8.DecodeRuneInString(text)
				if unicode.IsLetter(firstc) || unicode.IsDigit(firstc) {
					r.WriteByte(lex.ItemSpace)
				}
			}
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("code", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.EscapeHTML(node.Tokens)
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</code>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.spanNodeAttrs(node.Parent, &attrs)
		r.Tag("strong", attrs, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "bq")
		r.Tag("div", attrs, false)
	} else {
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		level := headingLevel[node.HeadingLevel : node.HeadingLevel+1]
		attrs = append(attrs, []string{"data-subtype", "h" + level})
		r.blockNodeAttrs(node, &attrs, "h"+level)
		r.Tag("div", attrs, false)
		attrs = [][]string{}
		r.contenteditable(&attrs)
		r.spellcheck(&attrs)
		r.Tag("div", attrs, false)
	} else {
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		switch node.ListData.Typ {
		case 0:
			attrs = append(attrs, []string{"data-subtype", "u"})
		case 1:
			attrs = append(attrs, []string{"data-subtype", "o"})
		case 3:
			attrs = append(attrs, []string{"data-subtype", "t"})
		}
		r.blockNodeAttrs(node, &attrs, "list")
		r.Tag("div", attrs, false)
	} else {
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		class := "li"
		var attrs [][]string
		switch node.ListData.Typ {
		case 0:
			attrs = append(attrs, []string{"data-marker", "*"})
			attrs = append(attrs, []string{"data-subtype", "u"})
		case 1:
			attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.ListData.Num) + "."})
			attrs = append(attrs, []string{"data-subtype", "o"})
		case 3:
			attrs = append(attrs, []string{"data-marker", "*"})
			attrs = append(attrs, []string{"data-subtype", "t"})
			if node.FirstChild.TaskListItemChecked {
				class += " protyle-task--done"
			}
		}
		r.blockNodeAttrs(node, &attrs, class)
		r.Tag("div", attrs, false)

		if 0 == node.ListData.Typ {
			attr := [][]string{{"class", "protyle-action"}, {"draggable", "true"}}
			r.Tag("div", attr, false)
			r.WriteString("<svg><use xlink:href=\"#iconDot\"></use></svg>")
			r.Tag("/div", nil, false)
		} else if 1 == node.ListData.Typ {
			attr := [][]string{{"class", "protyle-action protyle-action--order"}, {"contenteditable", "false"}, {"draggable", "true"}}
			r.Tag("div", attr, false)
			r.WriteString(strconv.Itoa(node.ListData.Num) + ".")
			r.Tag("/div", nil, false)
		}
	} else {
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if node.TaskListItemChecked {
			r.WriteString("<div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconCheck\"></use></svg></div>")
		} else {
			r.WriteString("<div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div>")
		}
		if nil == node.Next {
			node.InsertAfter(&ast.Node{ID: ast.NewNodeID(), Type: ast.NodeParagraph})
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "hr")
		r.Tag("div", attrs, false)
		r.Tag("div", nil, false)
	} else {
		r.Tag("/div", nil, false)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	return r.renderBr(node, entering)
}

func (r *BlockRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemNewline)
		if nil != node.Previous && (ast.NodeStrong == node.Previous.Type ||
			ast.NodeEmphasis == node.Previous.Type ||
			ast.NodeTag == node.Previous.Type ||
			ast.NodeStrikethrough == node.Previous.Type ||
			ast.NodeUnderline == node.Previous.Type ||
			ast.NodeKramdownSpanIAL == node.Previous.Type) &&
			nil != node.Next && bytes.Equal(util.CaretTokens, node.Next.Tokens) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) spanNodeAttrs(node *ast.Node, attrs *[][]string) {
	*attrs = append(*attrs, node.KramdownIAL...)
}

func (r *BlockRenderer) blockNodeAttrs(node *ast.Node, attrs *[][]string, class string) {
	r.nodeID(node, attrs)
	r.nodeIndex(node, attrs)
	r.nodeDataType(node, attrs)
	r.nodeClass(node, attrs, class)

	for _, ial := range node.KramdownIAL {
		if "id" == ial[0] {
			continue
		}
		*attrs = append(*attrs, []string{ial[0], ial[1]})
	}
}

func (r *BlockRenderer) nodeClass(node *ast.Node, attrs *[][]string, class string) {
	*attrs = append(*attrs, []string{"class", class})
}

func (r *BlockRenderer) nodeDataType(node *ast.Node, attrs *[][]string) {
	*attrs = append(*attrs, []string{"data-type", node.Type.String()})
}

func (r *BlockRenderer) nodeID(node *ast.Node, attrs *[][]string) {
	*attrs = append(*attrs, []string{"data-node-id", r.NodeID(node)})
}

func (r *BlockRenderer) nodeIndex(node *ast.Node, attrs *[][]string) {
	if nil == node.Parent || ast.NodeDocument != node.Parent.Type {
		return
	}

	*attrs = append(*attrs, []string{"data-node-index", strconv.Itoa(r.NodeIndex)})
	r.NodeIndex++
	return
}

func (r *BlockRenderer) spellcheck(attrs *[][]string) {
	*attrs = append(*attrs, []string{"spellcheck", "false"})
	return
}

func (r *BlockRenderer) contenteditable(attrs *[][]string) {
	*attrs = append(*attrs, []string{"contenteditable", strconv.FormatBool(r.Options.ProtyleContenteditable)})
	return
}

func (r *BlockRenderer) renderIAL(node *ast.Node) {
	attrs := [][]string{{"class", "protyle-attr"}, {"contenteditable", "false"}}
	r.Tag("div", attrs, false)

	if name := node.IALAttr("name"); "" != name {
		r.Tag("div", [][]string{{"class", "protyle-attr--name"}}, false)
		r.WriteString("<svg><use xlink:href=\"#iconN\"></use></svg>")
		r.WriteString(name)
		r.Tag("/div", nil, false)
	}

	if alias := node.IALAttr("alias"); "" != alias {
		r.Tag("div", [][]string{{"class", "protyle-attr--alias"}}, false)
		r.WriteString("<svg><use xlink:href=\"#iconA\"></use></svg>")
		r.WriteString(alias)
		r.Tag("/div", nil, false)
	}

	if memo := node.IALAttr("memo"); "" != memo {
		r.Tag("div", [][]string{{"class", "protyle-attr--memo b3-tooltips b3-tooltips__nw"}, {"aria-label", memo}}, false)
		r.WriteString("<svg><use xlink:href=\"#iconM\"></use></svg>")
		r.Tag("/div", nil, false)
	}

	if bookmark := node.IALAttr("bookmark"); "" != bookmark {
		r.Tag("div", [][]string{{"class", "protyle-attr--bookmark"}}, false)
		r.WriteString(bookmark)
		r.Tag("/div", nil, false)
	}

	if refCount := node.IALAttr("refcount"); "" != refCount {
		r.Tag("div", [][]string{{"class", "protyle-attr--refcount popover__block"}}, false)
		r.WriteString(refCount)
		r.Tag("/div", nil, false)
	}

	r.Tag("/div", nil, false)
}
