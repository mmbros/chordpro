package chordpro

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Debug(t *testing.T) {

	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "frontmatter+song1",
			src: `---    
title: una canzone
country: it
---
{t: titolo}
{st: sotto-titolo}{meta: album album}
[C]C'era una volta una [Am]gatta`,
			want: `FRONTMATTER
---    
title: una canzone
country: it
---

SONGS: count=1
  SONG 1
    META
      title: titolo
      subtitle: sotto-titolo
      album: album
    PAR 1 <Verse>
      LINE 1
        PAIR 1: [C] / C'era una volta una 
        PAIR 2: [Am] / gatta
`,
		},
		{
			name: "no-frontmatter+song-without-meta",
			src: `{soc}
la la la la la
{end_of_chorus}`,
			want: `SONGS: count=1
  SONG 1
    PAR 1 <Chorus>
      LINE 1
        PAIR 1:  / la la la la la
`,
		},
		{
			name: "par type",
			src: `{sov}[A]verse1[B]{eov}
{start_of_verse}
verse2
{end_of_verse}
{ns}
{start_of_bridge}
[C#m7]bridge
{end_of_bridge}
{soc}
[F]chorus 1
{eoc}
{new_song}
{start_of_chorus}
[G]chorus 2
{end_of_chorus}
{sot}
--- [A] tab 1
{eot}
{comment: comment  }
{start_of_tab}
--- [B] tab 2
{end_of_tab}
# do not print
{chorus}
`,
			want: `SONGS: count=3
  SONG 1
    PAR 1 <Verse>
      LINE 1
        PAIR 1: [A] / verse1
        PAIR 2: [B] / 
    PAR 2 <Verse>
      LINE 1
        PAIR 1:  / verse2
  SONG 2
    PAR 1 <Bridge>
      LINE 1
        PAIR 1: [C#m7] / bridge
    PAR 2 <Chorus>
      LINE 1
        PAIR 1: [F] / chorus 1
  SONG 3
    PAR 1 <Chorus>
      LINE 1
        PAIR 1: [G] / chorus 2
    PAR 2 <Tab>
      LINE 1
        PAIR 1:  / --- [A] tab 1
    PAR 3 <Comment>
      LINE 1
        PAIR 1:  / comment
    PAR 4 <Tab>
      LINE 1
        PAIR 1:  / --- [B] tab 2
    PAR 5 <ChorusRef>
`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			file := ParseText(tt.src)

			sb := strings.Builder{}
			df := NewDebugFormatter(&sb)
			df.WriteFile(file)
			got := sb.String()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("meta mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
