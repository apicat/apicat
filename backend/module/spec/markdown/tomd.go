package markdown

import (
	"bytes"
	"fmt"
	"github.com/apicat/apicat/backend/module/spec"
	"strings"
)

// ToMarkdown 将node文档转为markdown
func ToMarkdown(root *spec.Document) ([]byte, error) {
	r := &docRenderMarkdown{newline: true}
	for _, v := range root.Items {
		r.renderNode(v)
	}
	return r.buf.Bytes(), nil
}

type docRenderMarkdown struct {
	buf        bytes.Buffer
	listDepth  int
	linePrefix string
	lastText   string
	newline    bool
}

func (r *docRenderMarkdown) blockList(node *spec.DocNode) {
	r.listDepth++
	ordered := node.Type == "ordered_list"
	indent := strings.Repeat(" ", (r.listDepth-1)*3)
	for k, v := range node.Content {
		if ordered {
			r.content(fmt.Sprintf("%s%d. ", indent, k+1), v)
		} else {
			r.content(indent+"- ", v)
		}
	}
	r.listDepth--
	if r.listDepth == 0 {
		r.startLine()
		r.endLine()
	}
}

func (r *docRenderMarkdown) hr(flag bool) {
	r.buf.WriteString("\n")
	if flag {
		r.buf.WriteString("---")
		r.buf.WriteString("\n")
	}
	r.newline = true
}

func (r *docRenderMarkdown) startLine() {
	if r.newline {
		r.buf.WriteString(r.linePrefix)
	}
	r.newline = false
}

func (r *docRenderMarkdown) endLine() {
	if !r.newline {
		r.newline = true
		r.buf.WriteString("\n")
	}
}

func isNumber(data []byte) bool {
	for _, b := range data {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

func escape(text string) string {
	return string(bytes.Replace([]byte(text), []byte(`\`), []byte(`\\`), -1))
}

func needsEscaping(text string, lastNormalText string) bool {
	switch text {
	case `\`,
		"`",
		"*",
		"_",
		"{", "}",
		"[", "]",
		"(", ")",
		"#",
		"+",
		"-":
		return true
	case "!":
		return false
	case ".":
		return isNumber([]byte(lastNormalText))
	case "<", ">":
		return true
	default:
		return false
	}
}

func (r *docRenderMarkdown) content(prefix string, node *spec.DocNode) {
	r.startLine()
	r.buf.WriteString(prefix)
	for _, v := range node.Content {
		r.renderNode(v)
	}
	r.endLine()
}

func (r *docRenderMarkdown) inline(s string, marks []*spec.DocNode) {
	if len(marks) == 0 {
		r.buf.WriteString(s)
		return
	}
	for _, v := range marks {
		a, b := r.renderMark(v)
		r.buf.WriteString(a)
		r.buf.WriteString(s)
		r.buf.WriteString(b)
	}
}

func (r *docRenderMarkdown) renderMark(mark *spec.DocNode) (prefix, suffix string) {
	switch mark.Type {
	case "strong":
		return "**", "**"
	case "em":
		return "--", "--"
	case "link":
		title := mark.LookupAttrString("title")
		if title != "" {
			title = ` "` + title + `"`
		}
		return "[", "](" + escape(mark.LookupAttrString("href")) + title + ")"
	case "code":
		return "`", "`"
	default:
		return "", ""
	}
}

func (r *docRenderMarkdown) renderNode(node *spec.DocNode) {
	switch node.Type {
	// inline
	case "text":
		s := node.Text
		if needsEscaping(s, r.lastText) {
			s = "\\" + s
		}
		r.lastText = node.Text
		r.inline(s, node.Mark)
	case "image":
		r.inline(fmt.Sprintf("[%s](%s)",
			escape(node.LookupAttrString("src")), escape(node.LookupAttrString("title"))),
			node.Mark)
		// block
	case "heading":
		i := int(node.LookupAttrNumber("level"))
		r.content(strings.Repeat("#", i)+" ", node)
	case "paragraph":
		r.content("", node)
	case "blockquote":
		r.linePrefix = "> "
		r.content("", node)
		r.linePrefix = ""
	case "bullet_list", "ordered_list":
		r.blockList(node)
	case "list_item":
		r.content("", node)
	case "code_block":
		r.content("```", node)
	case "horizontal_rule":
		r.hr(true)
	case "hard_break":
		r.hr(false)
	default:
		// unspport op
	}
}
