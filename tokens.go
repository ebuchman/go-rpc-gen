package main

import (
	"fmt"
)

func (t token) String() string {
	s := fmt.Sprintf("Line %-2d, Col %-2d \t %-6s \t", t.loc.line, t.loc.col, t.typ.String())
	switch t.typ {
	case tokenEOFTy:
		return s + "EOF"
	case tokenErrTy:
		return s + t.val
	}
	/*if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}*/
	return s + fmt.Sprintf("%q", t.val)
}

// token types
type tokenType int

// TODO: use go generate
func (t tokenType) String() string {
	switch t {
	case tokenErrTy:
		return "[Error]"
	case tokenEOFTy:
		return "[EOF]"
	case tokenLeftBracesTy:
		return "[LeftBraces]"
	case tokenRightBracesTy:
		return "[RightBraces]"
	case tokenStringTy:
		return "[String]"
	case tokenSpaceTy:
		return "[Space]"
	}
	return "[Unknown]"
}

// token types
const (
	tokenErrTy            tokenType = iota // error
	tokenEOFTy                             // end of file
	tokenLeftBracesTy                      // {{
	tokenRightBracesTy                     // }}
	tokenLeftCurlBraceTy                   // {
	tokenRightCurlBraceTy                  //}
	tokenStringTy                          // variable name, contents of quotes, comments
	tokenLeftBraceTy                       // (
	tokenRightBraceTy                      // )
	tokenSpaceTy
)

// tokens
var (
	tokenLeftBraces     = "{{"
	tokenRightBraces    = "}}"
	tokenLeftCurlBrace  = "{"
	tokenRightCurlBrace = "}"
	tokenLeftBrace      = "("
	tokenRightBrace     = ")"
	tokenSpace          = " "
	tokenChars          = "abcdefghijklmnopqrstuvwqyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-/_.*,\n\t:+-/=`'\"!%&|[]>< {()"
)
