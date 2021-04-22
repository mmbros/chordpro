package chordpro

import (
	"fmt"

	"github.com/mmbros/chordpro/internal/lexer"
)

const (
	commentBegin     = '#'
	chordBegin       = '['
	chordEnd         = ']'
	directiveBegin   = '{'
	directiveEnd     = '}'
	directiveNameSep = ":"
	directiveMetaSep = " "
)

const (
	tokenNewline lexer.TokenType = iota
	tokenComment
	tokenText
	tokenChord
	tokenDirective
)

func trimDelim(src string) string {
	if len(src) <= 1 {
		return src
	}
	return src[1 : len(src)-1]
}

func stateText(l *lexer.L) lexer.StateFunc {
	// end with '{', '[', '\n', '\r', EOF
	var r rune
	var newState lexer.StateFunc

loop:
	for {
		r = l.Next()
		switch r {
		case commentBegin:
			newState = stateComment
			break loop
		case directiveBegin:
			newState = stateDirective
			break loop
		case chordBegin:
			newState = stateChord
			break loop
		case '\n', '\r':
			newState = stateNewline
			break loop
		case lexer.EOFRune:
			newState = nil
			break loop
		}
	}
	l.Rewind()
	if l.Current() != "" {
		l.Emit(tokenText)
	}

	return newState
}

func stateComment(l *lexer.L) lexer.StateFunc {
	// must begin with '#'
	// end with '\n', '\r', EOF

	r := l.Next()
	if r != commentBegin {
		l.Error(fmt.Sprintf("CommentState Invalid token: expected %q, got %q", commentBegin, r))
		return nil
	}
	// l.Ignore()

	for {
		r = l.Next()
		switch r {
		case '\n', '\r', lexer.EOFRune:
			l.Rewind()
			l.Emit(tokenComment)
			return stateNewline
		}
	}
}

func stateChord(l *lexer.L) lexer.StateFunc {
	// must begin with '['
	// end with ']'
	// error if '\n', '\r', EOF

	r := l.Next()
	if r != chordBegin {
		l.Error(fmt.Sprintf("ChordState Invalid token: expected %q, got %q", chordBegin, r))
		return nil
	}
	// l.Ignore()

	for {
		r = l.Next()
		switch r {
		case chordEnd:
			// l.Rewind()
			l.Emit(tokenChord)
			// l.Next()
			// l.Ignore()
			return stateText

		case '\n', '\r', lexer.EOFRune:
			ch := string(r)
			if r == lexer.EOFRune {
				ch = "EOF"
			}
			l.Error(fmt.Sprintf("ChordState Invalid token: expected %q, got %q", chordEnd, ch))
			return nil
		}
	}

}

func stateDirective(l *lexer.L) lexer.StateFunc {
	// must begin with '{'
	// end with '{'
	// error if '\n', '\r', EOF

	r := l.Next()
	if r != directiveBegin {
		l.Error(fmt.Sprintf("DirectiveState Invalid token: expected %q, got %q", directiveBegin, r))
		return nil
	}
	// l.Ignore()

	for {
		r = l.Next()
		switch r {
		case directiveEnd:
			// l.Rewind()
			l.Emit(tokenDirective)
			// l.Next()
			// l.Ignore()
			return stateText

		// case '\n', '\r', lexer.EOFRune:
		case lexer.EOFRune:
			ch := string(r)
			if r == lexer.EOFRune {
				ch = "EOF"
			}
			l.Error(fmt.Sprintf("DirectiveState Invalid token: expected %q, got %q", directiveEnd, ch))
			return nil
		}
	}

}

func stateNewline(l *lexer.L) lexer.StateFunc {
	// is '\n\r', '\n' or '\r'

	// \n   = CR (Carriage Return) // Used as a new line character in Unix
	// \r   = LF (Line Feed)       // Used as a new line character in Mac OS
	// \n\r = CR + LF              // Used as a new line character in Windows
	// (char)13 = \n = CR          // Same as \n

	var r rune
	precCR := false
	for {
		r = l.Next()
		switch r {
		case lexer.EOFRune:
			if precCR {
				l.Rewind()
				l.Emit(tokenNewline)
			}
			return nil
		case '\r':
			precCR = false
			l.Emit(tokenNewline)
		case '\n':
			if precCR {
				l.Rewind()
				l.Emit(tokenNewline)
				l.Next()
			}
			precCR = true
		default:
			if precCR {
				l.Rewind()
				l.Emit(tokenNewline)
			}
			return stateText
		}
	}
}
