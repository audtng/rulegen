package main

func ParseQuery(source *Source) (*QueryDocument, error) {
	p := parser{
		lexer:         lexer.New(source),
		maxTokenLimit: 0, // 0 means unlimited
	}
	return p.parseQueryDocument(), p.err
}

func ParseQueryWithTokenLimit(source *Source, maxTokenLimit int) (*QueryDocument, error) {
	p := parser{
		lexer:         lexer.New(source),
		maxTokenLimit: maxTokenLimit,
	}
	return p.parseQueryDocument(), p.err
}
