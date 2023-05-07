package markdown

import (
	"io"

	"github.com/apicat/apicat/common/spec"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// ToDocment 将markdown转为node文档
func ToDocment(md []byte) *spec.Document {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	renderer := &markdownRenderDoc{
		tree: make([]*spec.DocNode, 0),
	}
	markdown.Render(doc, renderer)
	return &spec.Document{Items: renderer.tree}
}

type markdownRenderDoc struct {
	tree  []*spec.DocNode
	state []*spec.DocNode
}

func (r *markdownRenderDoc) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.Text:
		r.inline("text", node.Literal)
	case *ast.Image:
		r.inline("image", nil, attr("src", string(node.Destination)))
	case *ast.Heading:
		r.block("heading", entering, attr("level", node.Level))
	case *ast.BlockQuote:
		r.block("blockquote", entering)
	case *ast.Paragraph:
		r.block("paragraph", entering)
	case *ast.List:
		switch node.ListFlags {
		case ast.ListTypeOrdered:
			r.block("ordered_list", entering)
		case ast.ListTypeTerm:
			r.block("bullet_list", entering)
		}
	case *ast.ListItem:
		r.block("list_item", entering)
	case *ast.CodeBlock:
		r.block("code_block", entering, attr("language", string(node.Info)))
	case *ast.Code:
		r.mark("code", entering)
	case *ast.Strong:
		r.mark("strong", entering)
	case *ast.Emph:
		r.mark("em", entering)
	case *ast.Link:
		r.mark("link",
			entering,
			attr("href", string(node.Destination)),
			attr("title", string(node.Title)),
		)
	default:
		// not support
	}
	return ast.GoToNext
}

func (r *markdownRenderDoc) RenderHeader(w io.Writer, ast ast.Node) {}
func (r *markdownRenderDoc) RenderFooter(w io.Writer, ast ast.Node) {}

func (r *markdownRenderDoc) statePop() {
	if n := len(r.state); n > 0 {
		r.state = r.state[:n-1]
	}
}

func (r *markdownRenderDoc) statePush(n *spec.DocNode) {
	last := r.stateNode()
	if last == nil {
		r.state = []*spec.DocNode{n}
		r.tree = append(r.tree, n)
	} else {
		last.Content = append(last.Content, n)
		r.state = append(r.state, n)
	}
}

func (r *markdownRenderDoc) stateNode() *spec.DocNode {
	n := len(r.state)
	if n == 0 {
		return nil
	}
	return r.state[n-1]
}

func (r *markdownRenderDoc) block(name string, entering bool, attrs ...attribute) {
	if entering {
		r.statePush(
			&spec.DocNode{
				Type:  name,
				Attrs: megredAttrs(attrs),
			},
		)
	} else {
		r.statePop()
	}
}

func (r *markdownRenderDoc) inline(name string, text []byte, attrs ...attribute) {
	n := &spec.DocNode{
		Type:  name,
		Attrs: megredAttrs(attrs),
	}
	if text != nil {
		n.Text = string(text)
	}
	cur := r.stateNode()
	cur.Content = append(cur.Content, n)
}

func (r *markdownRenderDoc) mark(name string, entering bool, attrs ...attribute) {
	if entering {
		cur := r.stateNode()
		cur.Mark = append(cur.Mark,
			&spec.DocNode{
				Type:  name,
				Attrs: megredAttrs(attrs),
			},
		)
	}
}

type attribute struct {
	Key   string
	Value any
}

func attr(k string, v any) attribute {
	return attribute{
		Key:   k,
		Value: v,
	}
}

func megredAttrs(attrs []attribute) map[string]any {
	if len(attrs) == 0 {
		return nil
	}
	mp := make(map[string]any, len(attrs))
	for _, v := range attrs {
		mp[v.Key] = v.Value
	}
	return mp
}
