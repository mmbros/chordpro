package chordpro

import (
	"fmt"
	"io"
	"strings"
)

// chord-sheet
//   paragraph
//     .tablatue
//     .comment
//     .verse
//     row
//       column
//         chord
//         lyrics

const (
	clsSong      = "chord-sheet"
	clsParagraph = "paragraph"
	clsLine      = "row"
	clsPair      = "column"
	clsChord     = "chord"
	clsLyric     = "lyrics"
	clsError     = "error"
)

const (
	tagSong         = "div"
	tagParagraph    = "div"
	tagParagraphPre = "pre"
	tagLine         = "div"
	tagPair         = "span"
	tagChord        = "u"
	tagLyric        = "i"
	tagError        = "div"
)

type HtmlDivFormatter struct {
	w io.Writer
}

func NewHtmlDivFormatter(w io.Writer) HtmlDivFormatter {
	return HtmlDivFormatter{
		w: w,
	}
}

func (f HtmlDivFormatter) appendTagOpen(tag, className string, newline bool) {
	fmt.Fprint(f.w, "<", tag)
	if className != "" {
		fmt.Fprint(f.w, ` class="`, className, `"`)
	}
	fmt.Fprint(f.w, ">")
	if newline {
		fmt.Fprintln(f.w)
	}
}
func (f HtmlDivFormatter) appendTagClose(tag string, newline bool) {
	fmt.Fprint(f.w, `</`, tag, `>`)
	if newline {
		fmt.Fprintln(f.w)
	}
}

func (f HtmlDivFormatter) appendParPre(className string, p *Paragraph) {
	f.appendTagOpen(tagParagraphPre, className, true)
	for _, lin := range p.Lines {
		for _, pair := range lin.Pairs {
			fmt.Fprint(f.w, pair.Lyric)
		}
		fmt.Fprintln(f.w)
	}
	f.appendTagClose(tagParagraphPre, true)
}

func (f HtmlDivFormatter) appendChordLyric(className string, pair, prec *ChordLyricPair) {

	txt := strings.TrimSpace(pair.Lyric)

	// handle newline
	// if the pair is in the middle of a word,
	// doesn't print a newline to keep the word together.
	// see: https://css-tricks.com/fighting-the-space-between-inline-block-elements
	if prec != nil {
		if strings.HasSuffix(prec.Lyric, " ") || (txt == "") {
			fmt.Fprintln(f.w)
		}
	}

	f.appendTagOpen(tagPair, className, false)

	f.appendTagOpen(tagChord, clsChord, false)
	fmt.Fprint(f.w, trimDelim(pair.Chord))
	f.appendTagClose(tagChord, false)

	f.appendTagOpen(tagLyric, clsLyric, false)
	if txt == "" {
		// to have chord always on top of the column even where there is no lyric
		txt = "&nbsp;"
	}
	fmt.Fprint(f.w, txt)
	f.appendTagClose(tagLyric, false)

	// never print newline at the end
	f.appendTagClose(tagPair, false)
}

func (f HtmlDivFormatter) appendLine(className string, lin *Line) {
	f.appendTagOpen(tagLine, className, true)
	var prec *ChordLyricPair
	for _, pair := range lin.Pairs {
		f.appendChordLyric(clsPair, pair, prec)
		prec = pair
	}
	f.appendTagClose(tagLine, true)
}

func (f HtmlDivFormatter) appendParagraph(p *Paragraph) {

	className := []string{"verse", "comment", "tablature", "chorus", "chorusref", "bridge"}[p.ParagraphType]

	switch p.ParagraphType {
	case Tab, Comment:
		f.appendParPre(className, p)
	case ChorusRef:
		f.appendTagOpen(tagParagraph, className, false)
		fmt.Fprint(f.w, "Chorus")
		f.appendTagClose(tagParagraph, true)
	default:
		f.appendTagOpen(tagParagraph, className, true)
		for _, lin := range p.Lines {
			f.appendLine(clsLine, lin)
		}
		f.appendTagClose(tagParagraph, true)
	}

}

func (f HtmlDivFormatter) appendParagraphs(ps []*Paragraph) {
	for _, p := range ps {
		f.appendParagraph(p)
	}
}

// func (f HtmlDivFormatter) appendFrontMatter(song *Song) {
// 	var s string
// 	fmt.Fprintln(f.w, `---`)

// 	s = song.Title()
// 	fmt.Fprintf(f.w, "title: %q\n", s)

// 	if s = song.Artist(); s == "" {
// 		s = song.SubTitle()
// 	}
// 	if s != "" {
// 		fmt.Fprintf(f.w, "artist: %q\n", s)
// 	}

// 	if s = song.Album(); s != "" {
// 		fmt.Fprintf(f.w, "album: %q\n", s)
// 	}

// 	if s = song.Year(); s != "" {
// 		fmt.Fprintf(f.w, "year: %q\n", s)
// 	}

// 	fmt.Fprintln(f.w, `---`)
// }

func (f HtmlDivFormatter) FormatBody(s *Song) {
	// f.appendFrontMatter(s)

	f.appendTagOpen(tagSong, clsSong, true)
	f.appendParagraphs(s.Paragraphs)
	if s.Err != nil {
		f.appendTagOpen(tagError, clsError, false)
		fmt.Fprint(f.w, s.Err.Error())
		f.appendTagClose(tagSong, false)
	}
	f.appendTagClose(tagSong, false)
	fmt.Fprint(f.w, "<!-- /", clsSong, " -->\n")
}
