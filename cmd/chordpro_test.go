package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_parseFrontmatter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  frontmatterMode
		err   error
	}{
		{
			name:  "ok-NONE",
			input: "NONE",
			want:  modeFrontmatterNone,
		},
		{
			name:  "ok-OverWrite",
			input: "OverWrite",
			want:  modeFrontmatterOverwrite,
		},
		{
			name:  "ok-preserve",
			input: "preserve",
			want:  modeFrontmatterPreserve,
		},
		{
			name:  "err-xxx",
			input: "xxx",
			err:   ErrInvalidFrontmatter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := parseFrontmatter(tt.input)
			if tt.err != nil {
				if tt.err != err {
					t.Errorf("expected %q error, got %q error", tt.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error %q", err.Error())
					return
				}

				if got != tt.want {
					t.Errorf("expected %v, got %v", tt.want, got)
				}
			}
		})
	}
}

func Test_parseOverwrite(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  overwriteMode
		err   error
	}{
		{
			name:  "ok-ALL",
			input: "ALL",
			want:  modeOverwriteAll,
		},
		{
			name:  "ok-NOne",
			input: "NOne",
			want:  modeOverwriteNone,
		},
		{
			name:  "ok-old",
			input: "old",
			want:  modeOverwriteOld,
		},
		{
			name:  "err-xxx",
			input: "xxx",
			err:   ErrInvalidOverwrite,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := parseOverwrite(tt.input)
			if tt.err != nil {
				if tt.err != err {
					t.Errorf("expected %q error, got %q error", tt.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error %q", err.Error())
					return
				}

				if got != tt.want {
					t.Errorf("expected %v, got %v", tt.want, got)
				}
			}
		})
	}
}

func Test_Walk(t *testing.T) {

	src := "../test/"

	err := filepath.Walk(src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relpath, _ := filepath.Rel(src, path)

			// dstpath := filepath.Join(dst, relpath)

			if !info.IsDir() {
				ext := filepath.Ext(path)
				switch strings.ToLower(ext) {
				case ".cho", ".chopro", ".chordpro":
					t.Log(relpath)

				}

				// t.Log(relpath, "->", dstpath)
			}

			return nil
		})
	if err != nil {
		t.Error(err)
	}
	// t.Fail()

}

func Test_transform(t *testing.T) {
	// func transform(r io.Reader, w io.Writer, prefix string, songFrontmatter bool) error {

	tests := []struct {
		name        string
		input       string
		want        string
		prefix      string
		frontmatter bool
		err         error
	}{
		{
			name:  "err-0",
			input: "",
			err:   ErrZeroSongs,
		},
		{
			name:  "err-N",
			input: "[C]do {new_song} [D]re",
			err:   ErrMultipleSongs,
		},
		{
			name:  "ok",
			input: "[C]do",
			want: `<div class="chord-sheet">
<div class="verse">
<div class="row">
<span class="column"><u class="chord">C</u><i class="lyrics">do</i></span></div>
</div>
</div><!-- /chord-sheet -->
`,
		},
		{
			name:   "ok-prefix",
			input:  "[C]do",
			prefix: "---\ntitle: \"titolo\"\n---\n",
			want: `---
title: "titolo"
---
<div class="chord-sheet">
<div class="verse">
<div class="row">
<span class="column"><u class="chord">C</u><i class="lyrics">do</i></span></div>
</div>
</div><!-- /chord-sheet -->
`,
		},
		{
			name:        "ok-frontmatter",
			input:       "{t: titolo}[C]do",
			frontmatter: true,
			want: `---
title: "titolo"
---
<div class="chord-sheet">
<div class="verse">
<div class="row">
<span class="column"><u class="chord">C</u><i class="lyrics">do</i></span></div>
</div>
</div><!-- /chord-sheet -->
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := strings.NewReader(tt.input)
			w := &strings.Builder{}

			err := transform(r, w, tt.prefix, tt.frontmatter)
			if tt.err != nil {
				if tt.err != err {
					t.Errorf("expected %q error, got %q error", tt.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error %q", err.Error())
					return
				}

				got := w.String()
				if got != tt.want {
					t.Errorf("expected %v, got %v", tt.want, got)
				}
			}
		})
	}
}
