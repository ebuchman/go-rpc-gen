package main

import (
	"fmt"
	"log"
	"time"
)

type parseStateFunc func(p *parser) parseStateFunc

type parser struct {
	l         *lexer
	last      token
	peekCount int // 1 if we've peeked

	txt  []string // surrounding go code
	jobs []Job    // things to do
}

func (p *parser) results() ([]string, []Job) {
	return p.txt, p.jobs
}

type Job struct {
	ident string
	args  []string
}

func Parser(input string) *parser {
	l := Lex(input)
	p := &parser{
		l:    l,
		txt:  []string{},
		jobs: []Job{},
	}
	return p
}

func (p *parser) next() token {
	if p.peekCount == 1 {
		p.peekCount = 0
		return p.last

	}
	p.last = <-p.l.Chan()
	return p.last
}

func (p *parser) peek() token {
	if p.peekCount == 1 {
		return p.last
	}
	p.next()
	p.peekCount = 1
	return p.last
}

func (p *parser) backup() {
	p.peekCount = 1
}

func (p *parser) run() error {
	for state := parseStateStart; state != nil; state = state(p) {
	}
	if p.last.typ == tokenErrTy {
		// return  err
	}
	return nil
}

// return a parseStateFunc that prints the error and triggers exit (returns nil)
func (p *parser) Error(s string) parseStateFunc {
	return func(pp *parser) parseStateFunc {
		// TODO: print location too
		log.Println("Error:", s)
		return nil
	}

}

func parseStateStart(p *parser) parseStateFunc {
	t := p.next()
	// scan past spaces, new lines, and comments
	switch t.typ {
	case tokenErrTy, tokenEOFTy:
		return nil
	//case tokenSpaceTy:
	//return parseStateStart
	case tokenStringTy, tokenLeftBraceTy, tokenRightBraceTy, tokenSpaceTy,
		tokenLeftCurlBraceTy, tokenRightCurlBraceTy:
		// write the text into the buffer
		p.txt = append(p.txt, t.val)
		return parseStateStart
	case tokenLeftBracesTy:
		return parseStateExpr
	}

	return p.Error(fmt.Sprintf("Unknown expression while parsing: %s", t.val))
}

// An expr contains an identifier that indicates which registered go functions
// need to be pasted in. It may have arguments itself.
func parseStateExpr(p *parser) parseStateFunc {
	var t = p.next()

	job := &Job{}
	for ; t.typ != tokenRightBracesTy; t = p.next() {
		//fmt.Println("StateExpr:", t.val)
		time.Sleep(10 * time.Millisecond)
		switch t.typ {
		case tokenStringTy:
			job.ident = t.val
		case tokenLeftBraceTy:
			// this identifier takes arguments
			p.parseArgs(job)
		case tokenErrTy, tokenEOFTy:
			break
		}
	}
	p.jobs = append(p.jobs, *job)
	return parseStateStart
}

func (p *parser) parseArgs(j *Job) {
	// TODO
}
