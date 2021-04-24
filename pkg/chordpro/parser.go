package chordpro

import (
	"fmt"
	"os"
	"strings"

	"github.com/mmbros/chordpro/internal/lexer"
)

type cursor struct {
	songs Songs
	song  *Song
	par   *Paragraph
	line  *Line
	pair  *ChordLyricPair

	onlyText bool
}

func (c *cursor) newSong() *Song {
	if c.songs == nil {
		c.songs = Songs{}
	} else {
		c.closeParagraph()
	}
	c.song = new(Song)
	c.song.meta = metaItems{}
	c.songs = append(c.songs, c.song)
	return c.song
}

func (c *cursor) newParagraph() *Paragraph {
	if c.song == nil {
		c.newSong()
	} else {
		c.closeLine()
	}
	c.par = new(Paragraph)
	c.song.Paragraphs = append(c.song.Paragraphs, c.par)

	return c.par
}

func (c *cursor) newLine() *Line {
	if c.par == nil {
		c.newParagraph()
	}
	c.line = new(Line)
	c.par.Lines = append(c.par.Lines, c.line)

	return c.line
}

func (c *cursor) newPair() *ChordLyricPair {
	if c.line == nil {
		c.newLine()
	}
	c.pair = new(ChordLyricPair)
	c.line.Pairs = append(c.line.Pairs, c.pair)

	return c.pair
}

func (c *cursor) getSong() *Song {
	if c.song == nil {
		c.newSong()
	}
	return c.song
}

func (c *cursor) getParagraph() *Paragraph {
	if c.par == nil {
		c.newParagraph()
	}
	return c.par
}

// func (c *cursor) getLine() *Line {
// 	if c.line == nil {
// 		c.newLine()
// 	}
// 	return c.line
// }

func (c *cursor) getPair() *ChordLyricPair {
	if c.pair == nil {
		c.newPair()
	}
	return c.pair

}

func (c *cursor) closeParagraph() {
	c.par = nil
	c.closeLine()
}

func (c *cursor) closeLine() {
	c.pair = nil
	c.line = nil
}

func (c *cursor) parseDirective(src string) {

	var name, arg string

	s := trimDelim(src)

	v := strings.SplitN(s, directiveNameSep, 2)
	name = strings.ToLower(strings.TrimSpace(v[0]))
	if len(v) > 1 {
		arg = strings.TrimSpace(v[1])
	}

	if name == "meta" {
		v = strings.SplitN(arg, directiveMetaSep, 2)
		if len(v) == 0 {
			return
		}
		name = v[0]
		if len(v) > 1 {
			arg = v[1]
		} else {
			arg = ""
		}
	}

	fieldName := metaNone

	switch name {
	case "title", "t":
		fieldName = metaTitle
	case "sorttitle":
		fieldName = metaSortTitle
	case "subtitle", "st":
		fieldName = metaSubtitle
	case "artist":
		fieldName = metaArtist
	case "composer":
		fieldName = metaComposer
	case "lyricist":
		fieldName = metaLyricist
	case "copyright":
		fieldName = metaCopyright
	case "album":
		fieldName = metaAlbum
	case "year":
		fieldName = metaYear
	case "key":
		fieldName = metaKey
	case "time":
		fieldName = metaTime
	case "tempo":
		fieldName = metaTempo
	case "duration":
		fieldName = metaDuration
	case "capo":
		fieldName = metaCapo

	case "comment", "c":
		c.closeParagraph()
		c.newPair().Lyric = arg
		c.getParagraph().ParagraphType = Comment
		c.closeParagraph()

	case "new_song", "ns":
		c.newSong()

	case "sov", "start_of_verse":
		p := c.newParagraph()
		p.Label = arg
	case "eov", "end_of_verse":
		c.closeParagraph()

	case "sob", "start_of_bridge":
		p := c.newParagraph()
		p.Label = arg
		p.ParagraphType = Bridge
	case "eob", "end_of_bridge":
		c.closeParagraph()

	case "soc", "start_of_chorus":
		p := c.newParagraph()
		p.Label = arg
		p.ParagraphType = Chorus
	case "eoc", "end_of_chorus":
		c.closeParagraph()
	case "chorus":
		p := c.newParagraph()
		p.Label = arg
		p.ParagraphType = ChorusRef
		c.closeParagraph()

	case "sot", "start_of_tab":
		c.onlyText = true
		p := c.newParagraph()
		p.Label = arg
		p.ParagraphType = Tab
	case "eot", "end_of_tab":
		c.onlyText = false
		c.closeParagraph()
	}

	if fieldName != metaNone {
		// add new meta item
		c.getSong().meta.append(fieldName, arg)
	}

}

func ParseText(src string) Songs {

	var newlineCounter int
	cur := cursor{}

	l := lexer.New(src, stateText)
	l.ErrorHandler = func(l *lexer.L) {
		cur.getSong().Err = l.Err
		fmt.Fprintln(os.Stderr, l.Err)
	}
	l.StartSync()

	for {

		tok, done := l.NextToken()
		if done {
			break
		}

		if tok.Type == tokenNewline {
			newlineCounter++
		} else {
			newlineCounter = 0
		}

		switch tok.Type {
		case tokenChord:
			if cur.onlyText {
				p := cur.getPair()
				// p.Lyric += "[" + tok.Value + "]"
				p.Lyric += tok.Value
			} else {
				p := cur.newPair()
				p.Chord = tok.Value
			}
		case tokenText:
			p := cur.getPair()
			p.Lyric += tok.Value
		case tokenNewline:

			if cur.getParagraph().ParagraphType == Tab {
				if newlineCounter > 1 {
					cur.newLine()
				}
				cur.closeLine()
				continue
			}

			cur.closeLine()
			if newlineCounter == 2 {
				cur.closeParagraph()
			}
		case tokenDirective:
			cur.parseDirective(tok.Value)
		}

		//		fmt.Printf("%v: %q\n", tok.Type, tok.Value)

	}

	return cur.songs
}
