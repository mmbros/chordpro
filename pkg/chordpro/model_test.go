package chordpro

import (
	"reflect"
	"strings"
	"testing"
)

func Test_metaItems_append(t *testing.T) {
	tests := []struct {
		name string
		mis  metaItems
		mi   metaItem
	}{
		{
			name: "first",
			mis:  metaItems{},
			mi:   metaItem{metaTitle, "Title1"},
		},
		{
			name: "second",
			mis:  metaItems{&metaItem{metaTitle, "Title1"}},
			mi:   metaItem{metaTitle, "Title2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			len1 := len(tt.mis)
			tt.mis.append(tt.mi.name, tt.mi.value)
			len2 := len(tt.mis)
			if len2 != len1+1 {
				t.Errorf("append: len: expected %d, got %d", len1+1, len2)
			}
		})
	}
}

func Test_metaItems_byFieldName(t *testing.T) {
	tests := []struct {
		name      string
		fieldName metaFieldName
		mis       metaItems
		want      []string
	}{
		{
			name:      "zero title",
			fieldName: metaTitle,
			mis:       metaItems{},
			want:      []string{},
		},
		{
			name:      "one title",
			fieldName: metaTitle,
			mis:       metaItems{&metaItem{metaTitle, "Title1"}},
			want:      []string{"Title1"},
		},
		{
			name:      "two titles",
			fieldName: metaTitle,
			mis:       metaItems{&metaItem{metaTitle, "Title1"}, &metaItem{metaTitle, "Title2"}},
			want:      []string{"Title1", "Title2"},
		},
		{
			name:      "two titles and one artist",
			fieldName: metaTitle,
			mis:       metaItems{&metaItem{metaTitle, "Title1"}, &metaItem{metaArtist, "Artist1"}, &metaItem{metaTitle, "Title2"}},
			want:      []string{"Title1", "Title2"},
		},
		{
			name:      "one artist and two titles",
			fieldName: metaArtist,
			mis:       metaItems{&metaItem{metaTitle, "Title1"}, &metaItem{metaTitle, "Title2"}, &metaItem{metaArtist, "Artist1"}},
			want:      []string{"Artist1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mis.byFieldName(tt.fieldName); !reflect.DeepEqual(got, tt.want) {
				if (len(tt.want) == 0) && (len(got) == 0) {
					return
				}
				t.Errorf("metaItems.byFieldName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParagraphType_String(t *testing.T) {
	tests := []struct {
		name string
		pt   ParagraphType
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pt.String(); got != tt.want {
				t.Errorf("ParagraphType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChordLyricPair_toString(t *testing.T) {
	type fields struct {
		Chord string
		Lyric string
	}
	type args struct {
		sb   *strings.Builder
		i    int
		spad string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ChordLyricPair{
				Chord: tt.fields.Chord,
				Lyric: tt.fields.Lyric,
			}
			p.toString(tt.args.sb, tt.args.i, tt.args.spad)
		})
	}
}

func TestLine_toString(t *testing.T) {
	type fields struct {
		Pairs []*ChordLyricPair
	}
	type args struct {
		sb   *strings.Builder
		i    int
		spad string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Line{
				Pairs: tt.fields.Pairs,
			}
			l.toString(tt.args.sb, tt.args.i, tt.args.spad)
		})
	}
}

func TestParagraph_toString(t *testing.T) {
	type fields struct {
		ParagraphType ParagraphType
		Label         string
		Lines         []*Line
	}
	type args struct {
		sb   *strings.Builder
		i    int
		spad string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Paragraph{
				ParagraphType: tt.fields.ParagraphType,
				Label:         tt.fields.Label,
				Lines:         tt.fields.Lines,
			}
			p.toString(tt.args.sb, tt.args.i, tt.args.spad)
		})
	}
}

func TestSong_toString(t *testing.T) {
	type fields struct {
		meta       metaItems
		Paragraphs []*Paragraph
	}
	type args struct {
		sb   *strings.Builder
		i    int
		spad string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Song{
				meta:       tt.fields.meta,
				Paragraphs: tt.fields.Paragraphs,
			}
			s.toString(tt.args.sb, tt.args.i, tt.args.spad)
		})
	}
}

func TestSongs_String(t *testing.T) {
	tests := []struct {
		name string
		ss   Songs
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ss.String(); got != tt.want {
				t.Errorf("Songs.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
