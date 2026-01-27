package markdown

import (
	"fmt"
	"maps"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func GetPaths(md string) []string {
	gd := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	var paths []string
	root := gd.Parser().Parse(text.NewReader([]byte(md)))
	ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		img, ok := node.(*ast.Image)
		if !ok {
			return ast.WalkContinue, nil
		}
		dest := string(img.Destination)
		if strings.Contains(dest, ".gif") {
			return ast.WalkContinue, nil
		}
		paths = append(paths, dest)

		return ast.WalkContinue, nil

	})
	return paths
}

func TransformMd(md string, imgMap map[string]string) string {
	for key, value := range maps.All(imgMap) {

		md = strings.ReplaceAll(md, key, fmt.Sprintf("\n%s\n", value))
	}
	return md
}
