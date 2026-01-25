package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"unicode"
)

const (
	EQUAL int = iota
	WORD
)

type token struct {
	tokenType int16
	value     string
}

func Lexer(filepath string) []token {
	if path.IsAbs(filepath) {
		filepath = path.Clean(filepath)
	}
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var tokens []token

	var buff []rune

	fileContent, e := io.ReadAll(file)
	if e != nil {
		panic(e)
	}
	charRunes := bytes.Runes(fileContent)

	for _, char := range charRunes {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' {
			buff = append(buff, char)
		} else if string(char) == "=" {
			tokens = append(tokens, token{value: string(buff), tokenType: int16(WORD)})
			buff = []rune{}
			tokens = append(tokens, token{value: "=", tokenType: int16(EQUAL)})
		} else if char == '\n' {
			tokens = append(tokens, token{value: string(buff), tokenType: int16(WORD)})
			buff = []rune{}
		}
	}
	fmt.Print(tokens)
	return tokens
}

type Parser struct {
	tokens []token
	index  int
}

func (l *Parser) previous() token {
	if l.index > 1 {
		return l.tokens[l.index-1]
	} else {
		return l.tokens[0]
	}
}

func (l *Parser) next() token {
	if l.index < len(l.tokens)-1 {
		return l.tokens[l.index+1]
	} else {
		return l.tokens[len(l.tokens)-1]
	}
}

func NewPraser(tokens []token) *Parser {
	return &Parser{tokens: tokens, index: 0}
}

func (l *Parser) Prase() ConfigObj {
	var obj ConfigObj
	for i, token := range l.tokens {
		l.index = i
		if token.tokenType == int16(EQUAL) {
			previous := l.previous()
			next := l.next()
			switch previous.value {
			case "PAT":
				obj.PAT = next.value
			case "Subtle":
				obj.Subtle = next.value
			case "Highlight":
				obj.Highlight = next.value
			case "Text":
				obj.Text = next.value
			case "Warning":
				obj.Warning = next.value
			case "Special":
				obj.Special = next.value
			}
		}
	}
	return obj
}
