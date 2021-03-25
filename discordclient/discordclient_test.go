package discordclient

import (
	"reflect"
	"strings"
	"testing"
)

func TestChunkString(t *testing.T) {
	var chunkStringTests = []struct {
		name string
		size int
		in   string
		out  []string
	}{
		{
			name: "empty",
			size: 3,
			in:   "",
			out:  []string{""},
		},
		{
			name: "one normal brake and one moved",
			size: 3,
			in:   "a\naaaa",
			out:  []string{"a", "aaa", "a"},
		},
		{
			name: "long",
			size: 3,
			in:   "1\n23456\n78\n\n9",
			out:  []string{"1", "234", "56", "78\n", "9"},
		},
		{
			name: "short newline break",
			size: 3,
			in:   "12\n3",
			out:  []string{"12", "3"},
		},
		{
			name: "newlines without break",
			size: 10,
			in:   "a\naa\naa",
			out:  []string{"a\naa\naa"},
		},
	}

	for _, tt := range chunkStringTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			chunked := chunkString(tt.in, tt.size)

			if !reflect.DeepEqual(tt.out, chunked) {
				t.Errorf("got %+v, want %+v", strings.Join(chunked, ","), strings.Join(tt.out, ","))
			}
		})
	}
}

func TestChunkStringLong(t *testing.T) {
	chunked := chunkString(strings.Repeat("aaaaaaaaaaaaaaaaaaa\naaaaaaaaaaaaaaaaaaaaa", 40000), 2000)
	if len(chunked) != 835 {
		t.Errorf("expected len 835 but was %d", len(chunked))
	}
}
