package chordpro

import (
	"fmt"
	"io"
)

const padstep = "  "

type DebugFormatter struct {
	w io.Writer
}

func NewDebugFormatter(w io.Writer) *DebugFormatter {
	return &DebugFormatter{w}
}

func (df *DebugFormatter) WriteFile(file *File) {
	df.WriteFrontmatter(file.Frontmatter)
	df.WriteSongs(file.Songs)
}

func (df *DebugFormatter) WriteFrontmatter(frontmatter string) {
	if frontmatter != "" {
		fmt.Fprintf(df.w, "FRONTMATTER\n%s\n", frontmatter)
	}
}

func (df *DebugFormatter) WriteSongs(songs Songs) {
	fmt.Fprintf(df.w, "SONGS: count=%d\n", len(songs))
	for i, s := range songs {
		df.WriteSong(s, i+1, padstep)
	}
}

func (df *DebugFormatter) WriteSong(song *Song, idx int, pad string) {
	fmt.Fprintf(df.w, "%sSONG %d\n", pad, idx)

	df.WriteMeta(song.meta, pad+padstep)

	for j, par := range song.Paragraphs {
		df.WriteParagraph(par, j+1, pad+padstep)
	}

}

func (df *DebugFormatter) WriteMeta(m metaItems, pad string) {
	if len(m) == 0 {
		return
	}

	fmt.Fprintf(df.w, "%sMETA\n", pad)
	for _, mi := range m {
		fmt.Fprintf(df.w, "%s%s: %s\n", pad+padstep, mi.name.String(), mi.value)
	}

}

func (df *DebugFormatter) WriteParagraph(p *Paragraph, idx int, pad string) {
	fmt.Fprintf(df.w, "%sPAR %d <%s>\n", pad, idx, p.ParagraphType.String())
	for j, lin := range p.Lines {
		df.WriteLine(lin, j+1, pad+padstep)
	}
}

func (df *DebugFormatter) WriteLine(l *Line, idx int, pad string) {
	fmt.Fprintf(df.w, "%sLINE %d\n", pad, idx)
	for j, pair := range l.Pairs {
		df.WriteChordLyricPair(pair, j+1, pad+padstep)
	}
}

func (df *DebugFormatter) WriteChordLyricPair(p *ChordLyricPair, idx int, pad string) {
	fmt.Fprintf(df.w, "%sPAIR %d: %s / %s\n", pad, idx, p.Chord, p.Lyric)
}

// func (p *Paragraph) toString(sb *strings.Builder, i int, spad string) {
// 	fmt.Fprintf(sb, "%sPAR #%d <%s>\n", spad, i, p.ParagraphType.String())
// 	for j, lin := range p.Lines {
// 		lin.toString(sb, j+1, spad+pad)
// 	}

// }
