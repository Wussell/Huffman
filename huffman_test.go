package main

import (
	"fmt"
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
			got := countChars(ex.s)
			for k, v := range got {
				if got[k] != ex.want[k] {
					log.Fatalf("got %v for %c. want %v\n", v, k, ex.want[k])
				}
			}
		})
	}
}

func TestmakeForest(t *testing.T) {
	examples := []struct {
		name   string
		counts map[rune]int
		want   []tree
	}{
		{
			name: "go go gophers",
			counts: map[rune]int{
				32:  2,
				101: 1,
				103: 3,
				104: 1,
				111: 3,
				112: 1,
				114: 1,
				115: 1,
			},
			want: []tree{
				{32, 2, nil, nil},
				{101, 1, nil, nil},
				{103, 3, nil, nil},
				{104, 1, nil, nil},
				{111, 3, nil, nil},
				{112, 1, nil, nil},
				{114, 1, nil, nil},
				{115, 1, nil, nil},
			},
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			forest := makeForest(ex.counts)
			fmt.Println(forest)
			gotMap := make(map[rune]int)
			for _, tree := range forest {
				gotMap[tree.char] = tree.weight
			}
			fmt.Println(gotMap)
			for i := 0; i < len(ex.want); i++ {
				count, ok := gotMap[ex.want[i].char]
				if ok != false {
					if count != ex.want[i].weight {
						log.Fatalf("got %v, want %v", count, ex.want[i].weight)
					}
				} else {
					log.Fatalf("character '%c' missing from created forest. value: %v", ex.want[i].char, ex.want[i].char)
				}
			}

			/*
				for i := 0; i < len(forest); i++ {
					count, ok := gotMap[forest[i].char]
					if ok != false {
						if count != forest[i].weight {
							log.Fatalf("got %v, want %v", )
						}
					}
					for c, count := range gotMap {

					}
					if forest[i] != ex.want[i] {
						log.Fatalf("got %v, want %v for entry %v\n", got[i], ex.want[i], i)
					}
				}
			*/
		})
	}
}
