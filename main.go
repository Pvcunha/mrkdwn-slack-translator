package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var mds = `# header

sample text.

[link](http://www.google.com)

## teste
`
var logger = slog.Default()

func renderParagraph(w io.Writer, p *ast.Paragraph, entering bool) {
	if entering {
		logger.Info("renderParagraph", "paragraph content", string(p.Content))
		io.WriteString(w, "<div>")
	} else {
		io.WriteString(w, "</div>\n")
	}
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if para, ok := node.(*ast.Paragraph); ok {
		renderParagraph(w, para, entering)
		return ast.GoToNext, true
	}

	return ast.GoToNext, false
}

func myNewRender() *html.Renderer {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: myRenderHook,
	}

	return html.NewRenderer(opts)
}

func main() {
	fmt.Println("aba")

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	md := []byte(mds)
	doc := p.Parse(md)

	fmt.Println("--- ast")
	ast.Print(os.Stdout, doc)
	fmt.Println()

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}

	renderer := html.NewRenderer(opts)

	html := markdown.Render(doc, renderer)

	fmt.Println("----- html")
	fmt.Printf("%s", html)

	newRender := myNewRender()
	customHtml := markdown.Render(doc, newRender)

	fmt.Println("---------- custom Render")
	fmt.Printf("%s", customHtml)
}
