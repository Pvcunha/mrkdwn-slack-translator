package slack

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/gomarkdown/markdown/ast"
)

var logger = slog.Default()

type Flags int

const (
	FlagsNone Flags = 0
	SkipHead  Flags = 1 << iota
)

type RendererOptions struct {
	Flags Flags
}

type Renderer struct {
	Opts RendererOptions
}

func NewRender(opts RendererOptions) *Renderer {
	return &Renderer{
		Opts: opts,
	}
}

func (r Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	placeComma := true
	switch node := node.(type) {
	case *ast.Heading:
		r.Heading(w, node, entering)
	case *ast.Text:
		r.Text(w, node)
	case *ast.HorizontalRule:
		r.HorizontalRule(w, node, entering)
	default:
		logger.Info("slack_RenderNode", "Unknown node %T", node)
		placeComma = false
	}

	next := ast.GetNextNode(node)
	if next != nil && placeComma {
		io.WriteString(w, ",")
	}
	return ast.GoToNext
}

func (r Renderer) RenderHeader(w io.Writer, ast ast.Node) {
	logger.Info("slack_RenderHeader")
	if r.Opts.Flags&SkipHead != 0 {
		return
	}

	io.WriteString(w, "{\"blocks\":[")
}

func (r Renderer) RenderFooter(w io.Writer, ast ast.Node) {
	if r.Opts.Flags&SkipHead != 0 {
		logger.Info("slack_RenderFooter", "skipHead", r.Opts.Flags&SkipHead)
		return
	}

	io.WriteString(w, "]}")
}

func (r *Renderer) Text(w io.Writer, text *ast.Text) {
	textType := "mrkdwn"
	parent := text.Parent
	isSection := true

	if _, ok := parent.(*ast.Heading); ok {
		textType = "plain_text"
		isSection = false
	}

	var txt string
	txt = `"text": { "type": "%s", "text": "%s" }`

	if isSection {
		txt = fmt.Sprintf(`{"type": "section",%s`, txt)

	}
	io.WriteString(w, fmt.Sprintf(txt, textType, text.Literal))

	if isSection {
		io.WriteString(w, "}")
	}
}

func (r *Renderer) Heading(w io.Writer, heading *ast.Heading, entering bool) {
	if entering {
		io.WriteString(w, "{\"type\": \"header\"")
	} else {
		io.WriteString(w, "}")
	}
}

func (r *Renderer) HorizontalRule(w io.Writer, hr *ast.HorizontalRule, entering bool) {
	logger.Info("slack_HorizontalRule")
	if entering {
		divider := `{"type": "divider"}`
		io.WriteString(w, divider)
	}
}
