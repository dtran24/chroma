package lexers

import (
	"testing"

	assert "github.com/alecthomas/assert/v2"
	"github.com/alecthomas/chroma/v2"
)

func TestYAMLInlineCommentsWithColon(t *testing.T) {
	source := `foo: "Foo" // default
bar: "Bar" // default: Bar
plain: value // comment: still comment
hash: value # comment: still comment
number: 123 // default: 456
boolean: true // default: false
sequence:
  - "Item" // default: Item
next: value
`

	assertYAMLTokens(t, source, []chroma.Token{
		{Type: chroma.NameTag, Value: "foo"},
		{Type: chroma.Punctuation, Value: ":"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.LiteralStringDouble, Value: `"Foo"`},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Comment, Value: "// default"},
		{Type: chroma.TextWhitespace, Value: "\n"},

		{Type: chroma.NameTag, Value: "bar"},
		{Type: chroma.Punctuation, Value: ":"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.LiteralStringDouble, Value: `"Bar"`},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Comment, Value: "// default: Bar"},
		{Type: chroma.TextWhitespace, Value: "\n"},

		{Type: chroma.NameTag, Value: "plain"},
		{Type: chroma.Punctuation, Value: ":"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Literal, Value: "value"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Comment, Value: "// comment: still comment"},
		{Type: chroma.TextWhitespace, Value: "\n"},

		{Type: chroma.NameTag, Value: "hash"},
		{Type: chroma.Punctuation, Value: ":"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Literal, Value: "value"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Comment, Value: "# comment: still comment"},
		{Type: chroma.TextWhitespace, Value: "\n"},

		{Type: chroma.NameTag, Value: "number"},
		{Type: chroma.Punctuation, Value: ":"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.LiteralNumber, Value: "123"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Comment, Value: "// default: 456"},
		{Type: chroma.TextWhitespace, Value: "\n"},

		{Type: chroma.NameTag, Value: "boolean"},
		{Type: chroma.Punctuation, Value: ":"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.KeywordConstant, Value: "true"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Comment, Value: "// default: false"},
		{Type: chroma.TextWhitespace, Value: "\n"},

		{Type: chroma.NameTag, Value: "sequence"},
		{Type: chroma.Punctuation, Value: ":"},
		{Type: chroma.TextWhitespace, Value: "\n  "},
		{Type: chroma.Text, Value: "- "},
		{Type: chroma.LiteralStringDouble, Value: `"Item"`},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Comment, Value: "// default: Item"},
		{Type: chroma.TextWhitespace, Value: "\n"},

		{Type: chroma.NameTag, Value: "next"},
		{Type: chroma.Punctuation, Value: ":"},
		{Type: chroma.TextWhitespace, Value: " "},
		{Type: chroma.Literal, Value: "value"},
		{Type: chroma.TextWhitespace, Value: "\n"},
	})
}

func TestYAMLInlineCommentBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []chroma.Token
	}{
		{
			name: "quoted slash slash stays string",
			source: `quoted_url: "http://example.test:8080/path // not a comment: still string"
`,
			expected: []chroma.Token{
				{Type: chroma.NameTag, Value: "quoted_url"},
				{Type: chroma.Punctuation, Value: ":"},
				{Type: chroma.TextWhitespace, Value: " "},
				{Type: chroma.LiteralStringDouble, Value: `"http://example.test:8080/path // not a comment: still string"`},
				{Type: chroma.TextWhitespace, Value: "\n"},
			},
		},
		{
			name: "single quoted slash slash stays string",
			source: `single_quoted_url: 'http://example.test:8080/path // not a comment: still string'
`,
			expected: []chroma.Token{
				{Type: chroma.NameTag, Value: "single_quoted_url"},
				{Type: chroma.Punctuation, Value: ":"},
				{Type: chroma.TextWhitespace, Value: " "},
				{Type: chroma.LiteralStringSingle, Value: `'http://example.test:8080/path // not a comment: still string'`},
				{Type: chroma.TextWhitespace, Value: "\n"},
			},
		},
		{
			name: "plain URL stays literal",
			source: `plain_url: http://example.test:8080/path
`,
			expected: []chroma.Token{
				{Type: chroma.NameTag, Value: "plain_url"},
				{Type: chroma.Punctuation, Value: ":"},
				{Type: chroma.TextWhitespace, Value: " "},
				{Type: chroma.Literal, Value: "http://example.test:8080/path"},
				{Type: chroma.TextWhitespace, Value: "\n"},
			},
		},
		{
			name: "slash slash without separation stays literal",
			source: `no_space_slash: value//not-comment
`,
			expected: []chroma.Token{
				{Type: chroma.NameTag, Value: "no_space_slash"},
				{Type: chroma.Punctuation, Value: ":"},
				{Type: chroma.TextWhitespace, Value: " "},
				{Type: chroma.Literal, Value: "value//not-comment"},
				{Type: chroma.TextWhitespace, Value: "\n"},
			},
		},
		{
			name: "hash without separation stays literal",
			source: `no_space_hash: value#not-comment # comment
`,
			expected: []chroma.Token{
				{Type: chroma.NameTag, Value: "no_space_hash"},
				{Type: chroma.Punctuation, Value: ":"},
				{Type: chroma.TextWhitespace, Value: " "},
				{Type: chroma.Literal, Value: "value#not-comment"},
				{Type: chroma.TextWhitespace, Value: " "},
				{Type: chroma.Comment, Value: "# comment"},
				{Type: chroma.TextWhitespace, Value: "\n"},
			},
		},
		{
			name: "plain URL followed by slash slash comment",
			source: `plain_url_comment: http://example.test:8080/path // comment: ok
`,
			expected: []chroma.Token{
				{Type: chroma.NameTag, Value: "plain_url_comment"},
				{Type: chroma.Punctuation, Value: ":"},
				{Type: chroma.TextWhitespace, Value: " "},
				{Type: chroma.Literal, Value: "http://example.test:8080/path"},
				{Type: chroma.TextWhitespace, Value: " "},
				{Type: chroma.Comment, Value: "// comment: ok"},
				{Type: chroma.TextWhitespace, Value: "\n"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertYAMLTokens(t, tt.source, tt.expected)
		})
	}
}

func assertYAMLTokens(t *testing.T, source string, expected []chroma.Token) {
	t.Helper()

	lexer := Get("yaml")
	assert.NotZero(t, lexer)

	tokens, err := chroma.Tokenise(chroma.Coalesce(lexer), nil, source)
	assert.NoError(t, err)
	assert.Equal(t, source, chroma.Stringify(tokens...))
	for _, token := range tokens {
		assert.NotEqual(t, chroma.Error, token.Type)
	}
	assert.Equal(t, expected, tokens)
}
