package lex

import (
	"strings"
)

type TokenType int

const (
	TokenText TokenType = iota
	TokenVariable
	TokenDelimiter
)

type Token struct {
	Type  TokenType
	Value string
}

func Tokenize(template string) []Token {
	var tokens []Token
	var currentToken strings.Builder

	remaining := template
	openDelimiter := "{{"
	closeDelimiter := "}}"
	openDelimiterLen := len(openDelimiter)
	closeDelimiterLen := len(closeDelimiter)

	for len(remaining) > 0 {
		openIndex := strings.Index(remaining, openDelimiter)
		closeIndex := strings.Index(remaining, closeDelimiter)

		if openIndex == -1 && closeIndex == -1 {
			currentToken.WriteString(remaining)
			break
		}

		if openIndex != -1 && (openIndex < closeIndex || closeIndex == -1) {
			if openIndex > 0 {

				currentToken.WriteString(remaining[:openIndex])
				tokens = append(tokens, Token{Type: TokenVariable, Value: currentToken.String()})
				currentToken.Reset()
			}

			tokens = append(tokens, Token{Type: TokenDelimiter, Value: openDelimiter})
			remaining = remaining[openIndex+openDelimiterLen:]
		} else if closeIndex != -1 && (closeIndex < openIndex || openIndex == -1) {

			if closeIndex > 0 {
				currentToken.WriteString(remaining[:closeIndex])
				tokens = append(tokens, Token{Type: TokenVariable, Value: currentToken.String()})
				currentToken.Reset()
			}

			tokens = append(tokens, Token{Type: TokenDelimiter, Value: closeDelimiter})
			remaining = remaining[closeIndex+closeDelimiterLen:]
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, Token{Type: TokenText, Value: currentToken.String()})
	}

	return tokens
}

func Interpolate(tokens []Token, data map[string]string) string {
	var result strings.Builder

	for _, token := range tokens {
		switch token.Type {
		case TokenText:
			result.WriteString(token.Value)
		case TokenVariable:
			if value, ok := data[token.Value]; ok {
				result.WriteString(value)
			}
		}
	}

	return result.String()
}
