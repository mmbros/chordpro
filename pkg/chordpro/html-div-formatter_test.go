package chordpro

import (
	"testing"
)

func Test_ParseText(t *testing.T) {
	src := `{t: titolo della canzone }{artist: artista}

[Am]Working Class [G]Hero is [D]somethin' to be[Am]
{comment: Tab from: http://www.guitartabs.cc/tabs/j/john_lennon/working_class_hero_crd_ver_3.html}
When they've tortured and scared you for twenty odd years.`
	f := ParseText(src)
	ss := f.Songs

	{
		// songs
		want := 1
		got := len(ss)
		if got != want {
			t.Errorf("Len(songs): want %d, got %d", want, got)
			t.FailNow()
		}
	}
	{
		// meta
		want := "titolo della canzone"
		got := ss[0].Title()
		if got != want {
			t.Errorf("song[0].Title: want %q, got %q", want, got)
		}
		want = "artista"
		got = ss[0].Artist()
		if got != want {
			t.Errorf("song[0].Artist: want %q, got %q", want, got)
		}
		want = ""
		got = ss[0].Album()
		if got != want {
			t.Errorf("song[0].Album: want %q, got %q", want, got)
		}
	}

} /*

func TestHtmlDivFormatter(t *testing.T) {
	src := `{t:Come Together}
{st:Beatles}
{sot}
The riff that they use in the song goes something like this
------------1-----
-------H----1--3--
------0--2--------
0--0--------------
{eot}

[Dm]Here come old flat top, He come grooving up slowly,

{c:Play riff}

[Dm]He wear no shoe shine, he got toe jam football
`
	file := ParseText(src)

	t.Log(file)

	var sb strings.Builder

	f := NewHtmlDivFormatter(&sb)
	f.FormatBody(file.Songs[0])
	t.Log(sb.String())

	t.FailNow()
}
*/
