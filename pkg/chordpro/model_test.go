package chordpro

import (
	"reflect"
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

func TestMetaFieldName_String(t *testing.T) {
	tests := []struct {
		mfn  metaFieldName
		want string
	}{
		{metaAlbum, "album"},
		{metaArtist, "artist"},
		{metaTitle, "title"},
		{metaSubtitle, "subtitle"},
		{metaFieldName(-1), "metaFieldName(-1)"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mfn.String(); got != tt.want {
				t.Errorf("metaFieldName.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
