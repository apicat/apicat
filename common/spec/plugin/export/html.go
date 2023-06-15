package export

import (
	_ "embed"
	"fmt"

	"github.com/apicat/apicat/common/spec"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed data/highlight/highlight.min.js
var highlightjs []byte

//go:embed data/highlight/mono-blue.min.css
var highlightcss []byte

func HTML(in *spec.Spec) ([]byte, error) {
	source, err := Markdown(in)
	if err != nil {
		return nil, err
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(source)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// return markdown.Render(doc, renderer), nil

	htmlbytes := markdown.Render(doc, renderer)

	return []byte(fmt.Sprintf(htmllayout, in.Info.Title, highlightcss, htmlbytes, highlightjs)), nil

	// md := goldmark.New(
	// 	goldmark.WithExtensions(extension.GFM),
	// 	goldmark.WithParserOptions(
	// 	// parser.WithAutoHeadingID(),
	// 	),
	// 	goldmark.WithRendererOptions(
	// 		html.WithHardWraps(),
	// 		html.WithXHTML(),
	// 		html.WithUnsafe(),
	// 	),
	// )
	// var buf bytes.Buffer
	// if err := md.Convert(source, &buf); err != nil {
	// 	return nil, err
	// }
	// return []byte(fmt.Sprintf(htmlheader, in.Info.Title, highlightcss, buf.String(), highlightjs)), nil
}

var htmllayout = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>%s</title>
	<style>
	body{font-family: -apple-system,BlinkMacSystemFont,"Segoe UI","Noto Sans",Helvetica,Arial,sans-serif,"Apple Color Emoji","Segoe UI Emoji";font-size: 16px;line-height: 1.5;word-wrap: break-word;}
	h1{padding-bottom: .3em;font-size: 3em;border-bottom: 1px solid #d0d7de;}
	h2{font-size: 2em;border-bottom: 1px solid #d0d7de;padding-bottom: .2em;}
	h3{font-size: 1.5em;}
	h4{font-size: 1.2em;}
	h2,h3,h4,h5,h6{font-weight: 500;}
	li+li{margin-top:.25em}
	pre>code{border-radius: 10px;line-height:1.45;}
	code{font-size:85%%;background-color: #eaeef3;color: #00193a;border-radius: 4px;padding:.2em .4em;font-family: ui-monospace,SFMono-Regular,SF Mono,Menlo,Consolas,Liberation Mono,monospace;}     
    hr{height: 1px;border:none;border-top:4px #eee solid;margin:4em 0}
    a{text-decoration: none;color: #0969da;}
    a:hover{text-decoration: underline;}
    table{border-collapse: collapse;width: 100%%; border:1px #d0d7de solid}
    th,td{border-bottom:1px #d0d7de solid;padding:6px 12px}
    tr:hover{background-color: rgb(238, 247, 250);}
    th{background-color:#f6f8fa;font-weight:500}
	blockquote{padding: 0 1em;color: #57606a;border-left: .25em solid #d0d7de;margin:10px 0}
	%s  
	</style>
</head>
<body>
	<div style="max-width:768px;margin:0 auto">%s</div>
	<script>%s</script>
	<script>hljs.highlightAll();</script>
</body>
</html>
`
