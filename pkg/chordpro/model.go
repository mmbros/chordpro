package chordpro

import (
	"fmt"
	"strings"
)

const pad = "  "

// Meta-data directives
// Each song can have meta-data associated, for example the song title.
// Meta-data are mostly used by programs that help organizing collections of ChordPro songs.

type metaFieldName int

const (
	metaNone metaFieldName = iota
	metaInvalid

	metaTitle
	metaSortTitle
	metaSubtitle
	metaArtist
	metaComposer
	metaLyricist
	metaCopyright
	metaAlbum
	metaYear
	metaKey
	metaTime
	metaTempo
	metaDuration
	metaCapo
)

type metaItem struct {
	name  metaFieldName
	value string
}

type metaItems []*metaItem

type Songs []*Song

type Song struct {
	meta       metaItems
	Paragraphs []*Paragraph
}

type ParagraphType int

const (
	Verse ParagraphType = iota
	Comment
	Tab
	Chorus
	ChorusRef
	Bridge
)

func (mis *metaItems) append(name metaFieldName, value string) {
	*mis = append(*mis, &metaItem{name, value})
}

// byFieldName method returns an array with all the values of the given field name.
func (mis metaItems) byFieldName(name metaFieldName) []string {
	var a []string

	for _, mi := range mis {
		if mi.name == name {
			a = append(a, mi.value)
		}
	}
	return a
}

// byFieldName1 method returns the first string value of the given field name.
func (mis metaItems) byFieldName1(name metaFieldName) string {
	for _, mi := range mis {
		if mi.name == name {
			return mi.value
		}
	}
	return ""
}

type Paragraph struct {
	ParagraphType ParagraphType
	Label         string
	Lines         []*Line
}

type Line struct {
	Pairs []*ChordLyricPair
}

type ChordLyricPair struct {
	Chord string
	Lyric string
}

func (pt ParagraphType) String() string {

	switch pt {
	case Verse:
		return "Verse"
	case Comment:
		return "Comment"
	case Tab:
		return "Tab"
	case Chorus:
		return "Chorus"
	case ChorusRef:
		return "ChorusRef"
	case Bridge:
		return "Bridge"
	default:
		return fmt.Sprintf("ParagraphType:%d", pt)
	}
}

func (p *ChordLyricPair) toString(sb *strings.Builder, i int, spad string) {
	fmt.Fprintf(sb, "%sPAIR #%d: %s / %s\n", spad, i, p.Chord, p.Lyric)
}

func (l *Line) toString(sb *strings.Builder, i int, spad string) {
	fmt.Fprintf(sb, "%sLINE #%d\n", spad, i)
	for j, pair := range l.Pairs {
		pair.toString(sb, j+1, spad+pad)
	}
}

func (p *Paragraph) toString(sb *strings.Builder, i int, spad string) {
	fmt.Fprintf(sb, "%sPAR #%d <%s>\n", spad, i, p.ParagraphType.String())
	for j, lin := range p.Lines {
		lin.toString(sb, j+1, spad+pad)
	}

}

func (s *Song) toString(sb *strings.Builder, i int, spad string) {
	fmt.Fprintf(sb, "%sSONG #%d\n", spad, i)

	for _, mi := range s.meta {
		fmt.Fprintf(sb, "%s> meta%d = %q\n", spad, mi.name, mi.value)
	}

	for j, par := range s.Paragraphs {
		par.toString(sb, j+1, spad+pad)
	}
}

func (ss Songs) String() string {

	var sb strings.Builder

	fmt.Fprintln(&sb) // newline

	for j, song := range ss {
		song.toString(&sb, j+1, pad)
	}
	return sb.String()
}

func (s *Song) Title() string {
	return s.meta.byFieldName1(metaTitle)
}

func (s *Song) SortTitle() string {
	return s.meta.byFieldName1(metaSortTitle)
}

func (s *Song) SubTitle() string {
	return s.meta.byFieldName1(metaSubtitle)
}

func (s *Song) Artist() string {
	return s.meta.byFieldName1(metaArtist)
}

func (s *Song) Album() string {
	return s.meta.byFieldName1(metaAlbum)
}

func (s *Song) Year() string {
	return s.meta.byFieldName1(metaYear)
}

// metaComposer
// metaLyricist
// metaCopyright
// metaKey
// metaTime
// metaTempo
// metaDuration
// metaCapo
