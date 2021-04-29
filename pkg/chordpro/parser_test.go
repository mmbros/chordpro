package chordpro

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_ParseMeta(t *testing.T) {

	tests := []struct {
		name string
		src  string
		want metaItems
	}{
		{
			name: "only one title",
			src:  `{t:    the title    }`,
			want: metaItems{
				{metaTitle, "the title"},
			},
		},
		{
			name: "only one title and one subtile",
			src: `{t: the_title}
			
			[A]la la [D]la la la

			{artist: the artist}
			xxx
			xxxx xxxx
			`,
			want: metaItems{
				{metaTitle, "the_title"},
				{metaArtist, "the artist"},
			},
		},
		{
			name: "some subtitles and a sort title",
			src: `{st: the_title_1}{subtitle: the_title_2}
{sorttitle: sort_title   }{st: the_title_3}`,
			want: metaItems{
				{metaSubtitle, "the_title_1"},
				{metaSubtitle, "the_title_2"},
				{metaSortTitle, "sort_title"},
				{metaSubtitle, "the_title_3"},
			},
		},
		{
			name: "all",
			src: `{t: title 1}{title: title 2}
{sorttitle: sort title}
{st: subtitle 1}{subtitle: subtitle 2}
{artist: artist}
{composer: composer}
{lyricist: lyricist}
{copyright: copyright}
{album: album}
{year: year}
{key: key}
{time: time}
{tempo: tempo}
{duration: duration}
{capo: capo}
{meta: capo value}
{meta: invalid}
{meta}
`,

			want: metaItems{
				{metaTitle, "title 1"},
				{metaTitle, "title 2"},
				{metaSortTitle, "sort title"},
				{metaSubtitle, "subtitle 1"},
				{metaSubtitle, "subtitle 2"},
				{metaArtist, "artist"},
				{metaComposer, "composer"},
				{metaLyricist, "lyricist"},
				{metaCopyright, "copyright"},
				{metaAlbum, "album"},
				{metaYear, "year"},
				{metaKey, "key"},
				{metaTime, "time"},
				{metaTempo, "tempo"},
				{metaDuration, "duration"},
				{metaCapo, "capo"},
				{metaCapo, "value"},
				{metaInvalid, "invalid: "},
				{metaInvalid, ": "},
			},
		},
	}

	opt := cmp.Comparer(func(a, b *metaItem) bool {
		return a.name == b.name && a.value == b.value
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := ParseText(tt.src)
			got := file.Songs[0].meta

			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("meta mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

func Test_ParseFrontmatter(t *testing.T) {

	tests := []struct {
		name            string
		src             string
		wantFrontmatter string
	}{
		{
			name: "frontmatter+song",
			src: `	

---    
title: La gatta
country: it
---	 	
{t: titolo}
[C]C'era una volta una [Am]gatta`,
			wantFrontmatter: `---    
title: La gatta
country: it
---	 	
`,
		},
		{
			name: "only frontmatter",
			src: `+++
title = "La gatta"
+++
`,
			wantFrontmatter: `+++
title = "La gatta"
+++
`,
		},
		{
			name: "only song",
			src: `{t: titolo}
[C]C'era una volta una [Am]gatta`,
			wantFrontmatter: "",
		},
		{
			name: "wrong frontmatter",
			src: `------
title: "La gatta"
---
`,
			wantFrontmatter: "",
		},
		{
			name: "char before frontmatter",
			src: `x
---
title: "La gatta"
---
`,
			wantFrontmatter: "",
		},
		{
			name: "false frontmatter",
			src: `
--
title: "La gatta"
---
`,
			wantFrontmatter: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			file := ParseText(tt.src)

			if file.Frontmatter != tt.wantFrontmatter {
				t.Errorf("Frontmatter: expected %v, got %v", tt.wantFrontmatter, file.Frontmatter)
			}
		})
	}
}
