package chordpro

import (
	"reflect"
	"testing"

	"github.com/mmbros/chordpro/internal/lexer"
)

func Test_trimDelim(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"0123456789", "12345678"},
		{"01", ""},
		{"0", ""},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			if got := trimDelim(tt.src); got != tt.want {
				t.Errorf("trimDelim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateFrontpage(t *testing.T) {
	tests := []struct {
		name  string
		char3 rune
		src   string
		want  *lexer.Token
	}{
		{
			name:  "fm1",
			char3: '-',
			src: `---
key1: val1
key2: val2
---
[A] la la la
`,
			want: &lexer.Token{
				Type: tokenFrontmatter,
				Value: `---
key1: val1
key2: val2
---
`,
			},
		},
		{
			name:  "txt1",
			char3: '+',
			src: `++
key1: val1
+++`,
			want: &lexer.Token{
				Type:  tokenText,
				Value: `++`,
			},
		},
		{
			name:  "fm not closed",
			char3: '=',
			src: `===
key1: val1
==`,
			want: &lexer.Token{
				Type: tokenFrontmatter,
				Value: `===
key1: val1
==`,
			},
		},
		{
			name:  "fm with /r",
			char3: '-',
			src: `---  ` + string('\r') + `
key1: val1
---

{t: song}
`,
			want: &lexer.Token{
				Type: tokenFrontmatter,
				Value: `---  ` + string('\r') + `
key1: val1
---
`,
			},
		},
	}
	for _, tt := range tests {
		if tt.name != "fm with /r" {
			continue
		}

		t.Run(tt.name, func(t *testing.T) {

			stateFP := stateFrontpage(tt.char3)
			L := lexer.New(tt.src, stateFP)
			L.StartSync()

			got, _ := L.NextToken()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stateFrontpage Current: got %v, want %v", got, tt.want)
			}

		})
	}
}

// func Test_stateTrimInitialWhitespace(t *testing.T) {
// 	type args struct {
// 		l *lexer.L
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want lexer.StateFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := stateTrimInitialWhitespace(tt.args.l); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("stateTrimInitialWhitespace() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_stateText(t *testing.T) {
// 	type args struct {
// 		l *lexer.L
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want lexer.StateFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := stateText(tt.args.l); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("stateText() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_stateComment(t *testing.T) {
// 	type args struct {
// 		l *lexer.L
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want lexer.StateFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := stateComment(tt.args.l); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("stateComment() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_stateChord(t *testing.T) {
// 	type args struct {
// 		l *lexer.L
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want lexer.StateFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := stateChord(tt.args.l); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("stateChord() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_stateDirective(t *testing.T) {
// 	type args struct {
// 		l *lexer.L
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want lexer.StateFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := stateDirective(tt.args.l); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("stateDirective() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_stateNewline(t *testing.T) {
// 	type args struct {
// 		l *lexer.L
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want lexer.StateFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := stateNewline(tt.args.l); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("stateNewline() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
