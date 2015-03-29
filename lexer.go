package main

// This lexer is heavily inspired by Rob Pike's "Lexical Scanning in Go"

import (
	"fmt"
	"log"
	"strings"
)

type lexStateFunc func(*lexer) lexStateFunc

// a token
type token struct {
	typ tokenType
	val string

	loc location
}

// the lexer object
type lexer struct {
	input  string // input string to lex
	length int    // length of the input string
	pos    int    // current pos
	start  int    // start of current token

	line        int // current line number
	lastNewLine int // pos of last new line

	tokens chan token // channel to emit tokens over

	temp string // a place to hold eg. commands
}

// location for error reporting
type location struct {
	line int
	col  int
}

// Lex the input, returning the lexer
// Tokens can be fetched off the channel
func Lex(input string) *lexer {
	l := &lexer{
		input:  input,
		length: len(input),
		pos:    0,
		tokens: make(chan token, 2),
	}
	go l.run()
	return l
}

func (l *lexer) Error(s string) lexStateFunc {
	return func(l *lexer) lexStateFunc {
		// TODO: print location data too
		log.Println(s)
		return nil
	}
}

// Return the tokens channel
func (l *lexer) Chan() chan token {
	return l.tokens
}

// Run the lexer
func (l *lexer) run() {
	for state := lexStateStart; state != nil; state = state(l) {
	}
	l.emit(tokenEOFTy)
	close(l.tokens)
}

// Return next character in the string
// To hell with utf8 :p
func (l *lexer) next() string {
	if l.pos >= l.length {
		return ""
	}
	b := l.input[l.pos : l.pos+1]
	l.pos += 1
	return b
}

// backup a step
func (l *lexer) backup() {
	l.pos -= 1
}

// peek ahead a character without consuming
func (l *lexer) peek() string {
	s := l.next()
	l.backup()
	return s
}

// consume a token and push out on the channel
func (l *lexer) emit(ty tokenType) {
	l.tokens <- token{
		typ: ty,
		val: l.input[l.start:l.pos],
		loc: location{
			line: l.line,
			col:  l.pos - l.lastNewLine,
		},
	}
	l.start = l.pos
}

func (l *lexer) accept(options string) bool {
	if strings.Contains(options, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(options string) bool {
	i := 0
	s := l.next()
	for ; l.pos < l.length && strings.Contains(options, s); s = l.next() {
		i += 1
	}
	if strings.Contains(options, s) {
		// if we are at the end
		// the loop never runs
		return true
	} else if l.pos < l.length {
		l.backup()
	} else if s != "" {
		l.backup()
	}

	return i > 0
}

// Starting state
func lexStateStart(l *lexer) lexStateFunc {

	for {
		if strings.HasPrefix(l.input[l.pos:], tokenLeftBraces) {
			if l.pos > l.start {
				l.emit(tokenStringTy)
			}
			return lexStateLeftBraces // Next state.
		}
		if l.next() == "" {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(tokenStringTy)
	}
	l.emit(tokenEOFTy) // Useful to make EOF a token.
	return nil         // Stop the run loop.

	// check the one character tokens
	t := l.next()
	switch t {
	case "":
		return nil
	case tokenLeftBrace:
		l.emit(tokenLeftBraceTy)
		return lexStateStart
		/*case tokenRightBrace:
			l.emit(tokenRightBraceTy)
			return lexStateStart
		case tokenRightCurlBrace:
			l.emit(tokenRightCurlBraceTy)
			return lexStateStart*/
	}
	l.backup()

	remains := l.input[l.pos:]

	// skip spaces
	if isSpace(l.peek()) {
		return lexStateSpace
	}

	// check for left braces
	if strings.HasPrefix(remains, tokenLeftBraces) {
		return lexStateLeftBraces
	}
	return lexStateString
}

func isSpace(s string) bool {
	return s == " " || s == "\t"
}

func lexStateExpr(l *lexer) lexStateFunc {
	s := l.next()
	s2 := l.peek()
	for ; s+s2 != tokenRightBraces; s, s2 = l.next(), l.peek() {
		// check for chars
		if strings.Contains(tokenChars, s) {
			l.backup()
			if !l.acceptRun(tokenChars) {
				return l.Error("Expected a string")
			}
			l.emit(tokenStringTy)
		} else if s == tokenRightCurlBrace {
			// this guy gets special treatment incase he's }}
			l.emit(tokenRightCurlBraceTy)
		} else {
			return l.Error(fmt.Sprintf("Invalid char: %s", s))
		}
	}
	l.pos += 1 // consume the second brace
	l.emit(tokenRightBracesTy)
	return lexStateStart
}

// Scan past spaces
func lexStateSpace(l *lexer) lexStateFunc {
	for s := l.next(); isSpace(s); s = l.next() {
	}
	l.backup()
	l.start = l.pos
	return lexStateStart
}

// On {{
func lexStateLeftBraces(l *lexer) lexStateFunc {
	l.pos += len(tokenLeftBraces)
	l.emit(tokenLeftBracesTy)
	return lexStateExpr
}

// On }}
func lexStateRightBraces(l *lexer) lexStateFunc {
	l.pos += len(tokenRightBraces)
	l.emit(tokenRightBracesTy)
	return lexStateStart
}

// a string
func lexStateString(l *lexer) lexStateFunc {
	s := l.peek()
	if !l.acceptRun(tokenChars) {
		return l.Error(fmt.Sprintf("Expected a string. Got: %s", s))
		l.backup()
	}
	l.emit(tokenStringTy)
	return lexStateStart
}
