package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Pvcunha/mrkdwn-slack-translator/pkg/slack"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

var mds = `# header
sample text
---
teste`
var logger = slog.Default()

func main() {

	dat, err := os.ReadFile("a.md")
	if err != nil {
		panic("error while reading file")
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	doc := p.Parse(dat)

	fmt.Println("--- ast")
	ast.Print(os.Stdout, doc)
	fmt.Println()

	fmt.Println("------------ slack")
	// slackFlags := slack.SkipHead
	slackRenderer := slack.NewRender(slack.RendererOptions{Flags: 0})
	slack := markdown.Render(doc, slackRenderer)
	fmt.Println(string(slack))
}
