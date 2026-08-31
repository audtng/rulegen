package main

func ParseQuery(source *Source) (*QueryDocument, error) {
	p := parser{
		lexer:         lexer.New(source),
		maxTokenLimit: 0, // 0 is the default value
	}
	return p.parseQueryDocument(), p.err
}
