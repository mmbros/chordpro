package cmd

import (
	"strings"
	"testing"

	"github.com/mmbros/chordpro/pkg/chordpro"
)

func Test_appendFrontMatter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "full",
			input: " {t: title} \n\n {st: artist} xxx {year: 2021}{album:   album1  } x x x x {album:  album2  }",
			want: `---
title: "title"
artist: "artist"
album: "album1"
year: "2021"
---
`,
		},
		{
			name:  "empty",
			input: "lyric",
			want: `---
title: ""
---
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// ss := ParseText(src)
			file := chordpro.ParseText(tt.input)
			songs := file.Songs
			builder := strings.Builder{}

			appendFrontMatter(&builder, songs[0])

			got := builder.String()

			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func Test_getFrontMatterFromReader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "ok-YAML-1",
			input: `---
campo1: valore1
campo2: valore2
---
linea 1
linea 2
linea N`,
			want: `---
campo1: valore1
campo2: valore2
---
`,
		},
		{
			name: "ok-YAML-2",
			input: `---
campo1: valore1
+++
---
`,
			want: `---
campo1: valore1
+++
---
`,
		},
		{
			name: "WARN-YAML-3(newline added)",
			input: `---
campo1: valore1
---`,
			want: `---
campo1: valore1
---
`, // a newline is added
		},
		{

			name: "ok-YAML-4(unclosed and skip empty lines)",
			input: `


---
campo1: valore1
`,
			want: `---
campo1: valore1
`,
		},
		{

			name: "ok-YAML-5(no front matter)",
			input: `

linea 1
`,
			want: ``,
		},
		{
			name: "ok-TOM-1",
			input: `
   
+++
campo1 = valore1
---
campo2 = valore2
+++
linea 1
linea 2
linea N`,
			want: `+++
campo1 = valore1
---
campo2 = valore2
+++
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got := getFrontMatterFromReader(r)

			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
