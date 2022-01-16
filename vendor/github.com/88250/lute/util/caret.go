// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package util

// Caret 插入符 \u2038。
const Caret = "‸"

// CaretNewline 插入符加换行。
const CaretNewline = Caret + "\n"

// CaretTokens 是插入符的字节数组。
var CaretTokens = []byte(Caret)

// CaretRune 是插入符的 Rune。
var CaretRune = []rune(Caret)[0]

// CaretNewlineTokens 插入符加换行字节数组。
var CaretNewlineTokens = []byte(CaretNewline)

// CaretReplacement 用于解析过程中临时替换。
const CaretReplacement = "caretreplacement"

// FrontEndCaret 前端插入符。
const FrontEndCaret = "<wbr>"

// FrontEndCaret 前端自动闭合插入符。
const FrontEndCaretSelfClose = "<wbr/>"