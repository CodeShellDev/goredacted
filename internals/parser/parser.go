package parser

type Parser struct {
	Name string
	Parse func(this Parser, content string) string
}

var registry []Parser

func Parse(content string) string {
	for _, p := range registry {
		content = p.Parse(p, content)
	}

	return content
}

func Register(parser Parser) {
	registry = append(registry, parser)
}