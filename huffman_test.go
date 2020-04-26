package huffman

import (
	"log"
	"testing"
)

func TestCharCount(t *testing.T) {
	examples := []struct {
		name string
		s    string
		want map[rune]int
	}{
		{
			name: "green eggs and ham",
			s:    "I do not like green eggs and ham. I do not like them Sam-I-Am.",
			want: map[rune]int{
				32:  13,
				45:  2,
				46:  2,
				65:  1,
				73:  3,
				83:  1,
				97:  3,
				100: 3,
				101: 6,
				103: 3,
				104: 2,
				105: 2,
				107: 2,
				108: 2,
				109: 4,
				110: 4,
				111: 4,
				114: 1,
				115: 1,
				116: 3,
			},
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			got := charCount(ex.s)
			for k, v := range got {
				if got[k] != ex.want[k] {
					log.Fatalf("got %v for %c. want %v\n", v, k, ex.want[k])
				}
			}
		})
	}
}
