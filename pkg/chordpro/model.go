package chordpro

import (
	"fmt"
)

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

type File struct {
	Frontmatter string
	Songs       Songs
}

type Songs []*Song

type Song struct {
	meta       metaItems
	Paragraphs []*Paragraph
	Err        error
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

func (mfn metaFieldName) String() string {
	a := []string{
		"none",
		"invalid",
		"title",
		"sortTitle",
		"subtitle",
		"artist",
		"composer",
		"lyricist",
		"copyright",
		"album",
		"year",
		"key",
		"time",
		"tempo",
		"duration",
		"capo",
	}
	if mfn >= metaNone && mfn <= metaCapo {
		return a[mfn]
	}
	return fmt.Sprintf("metaFieldName(%d)", mfn)
}

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

// func (mi *metaItem) String() string {
// 	return
// }

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
		return fmt.Sprintf("ParagraphType(%d)", pt)
	}
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
