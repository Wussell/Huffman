package main

import (
	"testing"
)

/*
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

func TestMakeForest(t *testing.T) {
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
		})
	}
}

func TestMakeTree(t *testing.T) {
	examples := []struct {
		name   string
		forest []tree
	}{
		{
			name: "streets are stone stars are not",
			forest: []tree{
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

		})
	}
}
*/
func TestCompressTree(t *testing.T) {
	examples := []struct {
		name string
		root *tree
	}{
		{
			name: "streets are stone stars are not",
			root: makeTree("streets are stone stars are not"),
		},
		{
			name: "go go gophers",
			root: makeTree("go go gophers"),
		},

		{
			name: "the quick brown fox jumped over the lazy dog",
			root: makeTree("the quick brown fox jumped over the lazy dog"),
		},
		{
			name: "how much wood could a woodchuck chuck if a woodchuck could chuck wood",
			root: makeTree("how much wood could a woodchuck chuck if a woodchuck could chuck wood"),
		},
		{
			name: "she sells sea shells by the sea shore",
			root: makeTree("she sells sea shells by the sea shore"),
		},
		{
			name: "peter piper picked a peck of pickled peppers",
			root: makeTree("peter piper picked a peck of pickled peppers"),
		},

		{
			name: "abÞ",
			root: makeTree("abÞ"),
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			var serialTree string
			serialTree = compressTree(ex.root, serialTree)
			compressedTree := compressedTreeToBits(serialTree)
			treeEnd := findTreeEnd(compressedTree)
			uncompressedNodes := uncompressTree(compressedTree[:treeEnd])
			mNodes := make(map[int]tree)
			nodes := traverse(ex.root)

			for _, n := range nodes {
				mNodes[n.id] = n
			}
			for _, uN := range uncompressedNodes {
				wantNode, ok := mNodes[uN.id]
				if ok == false {
					t.Fatalf("There is no node with id %v\n", uN.id)
				}
				gotNode := *uN
				if gotNode.c != wantNode.c {
					t.Fatalf("Mismatched characters. got %v, want %v\n", gotNode.c, wantNode.c)
				}
			}
			if got, want := len(uncompressedNodes), len(nodes); got != want {
				t.Fatalf("Missing nodes from uncompressed tree. got %d, want %d\n", got, want)
			}

			uRoot := findRoot(uncompressedNodes)
			if uRoot.id != ex.root.id {
				t.Fatalf("Wrong root chosen. got %v, want %v.\n", *uRoot, *ex.root)
			}
			orderedUNodes := traverse(uRoot)
			for i, n := range nodes {
				if orderedUNodes[i].id != n.id {
					t.Fatalf("Nodes out of order. got %v, want %v at index %v.\n", orderedUNodes[i], nodes[i], i)
				}
			}
			var p string
			pathTable := paths(ex.root, p)
			p = ""
			uPathTable := paths(uRoot, p)
			for c, path := range pathTable {
				uPath, ok := uPathTable[c]
				if ok == false {
					t.Fatalf("character %c does not exist in decompressed tree", c)
				}
				if uPath != path {
					t.Fatalf("got %s, want %s for character %c. \n", uPath, path, c)
				}
			}
		})
	}
}

func TestCompress(t *testing.T) {
	examples := []struct {
		name string
		root *tree
	}{
		{
			name: "streets are stone stars are not" + "Þ",
			root: makeTree("streets are stone stars are not" + "Þ"),
		},
		{
			name: "go go gophers" + "Þ",
			root: makeTree("go go gophers" + "Þ"),
		},

		{
			name: "the quick brown fox jumped over the lazy dog" + "Þ",
			root: makeTree("the quick brown fox jumped over the lazy dog" + "Þ"),
		},
		{
			name: "how much wood could a woodchuck chuck if a woodchuck could chuck wood" + "Þ",
			root: makeTree("how much wood could a woodchuck chuck if a woodchuck could chuck wood" + "Þ"),
		},
		{
			name: "she sells sea shells by the sea shore" + "Þ",
			root: makeTree("she sells sea shells by the sea shore" + "Þ"),
		},
		{
			name: "peter piper picked a peck of pickled peppers" + "Þ",
			root: makeTree("peter piper picked a peck of pickled peppers" + "Þ"),
		},

		{
			name: "abÞ",
			root: makeTree("abÞ"),
		},
	}
	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			compressedTree := compressedTreeToBits(compressTree(ex.root, ""))
			compressedBits := stringToBits(ex.name, paths(ex.root, ""))
			compressedData := append(compressedTree, compressedBits...)

			treeEnd := findTreeEnd(compressedData)
			uNodes := uncompressTree(compressedData[:treeEnd])
			root := findRoot(uNodes)
			//	log.Printf("compressed tree: %x", compressedTree)
			//	for _, uNode := range uNodes {
			//		log.Printf("uncompressed tree nodes: %+v", uNode)
			//	}
			/*	log.Printf("compressed bits: %x", compressedBits)
				log.Printf("compressed tree end: %d", treeEnd)
				log.Printf("compressed root id: %d", root.id)
			*/
			reslicedData := compressedData[treeEnd+8:]
			for i, b := range reslicedData {
				if b != compressedBits[i] {
					t.Fatalf("compressed data resliced incorrectly at index %v", i)
				}
			}
			uString := unhuff(compressedData[treeEnd+8:], root)
			if uString != ex.name[:len(ex.name)-2] {
				t.Fatalf("strings don't match\ngot: %s \nwant: %s\n ", uString, ex.name[:len(ex.name)-2])
			}
		})
	}
}
