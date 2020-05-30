package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type tree struct {
	id int
	c  rune
	w  int
	l  *tree
	r  *tree
}

func countChars(content string) map[rune]int {
	count := make(map[rune]int)
	for _, v := range content {
		count[v]++
	}
	return count
}

func makeForest(charCount map[rune]int) []tree {
	var numChars int
	for range charCount {
		numChars++
	}
	//fmt.Println("number of characters: ", numChars)
	forest := make([]tree, numChars)
	i := 0
	for c, v := range charCount {
		forest[i].c = c
		forest[i].w = v
		//fmt.Println(forest[i])
		i++
	}
	return forest
}

func combineTrees(t1 tree, t2 tree) *tree {
	var t3 *tree = &tree{}
	t3.w = t1.w + t2.w
	t3.l = &t1
	t3.r = &t2
	return t3
}

func makeTree(forest []tree) *tree {
	var newTree *tree
	for len(forest) > 1 {
		sort.Slice(forest, func(i, j int) bool { return forest[i].w < forest[j].w })
		newTree = combineTrees(forest[0], forest[1])
		forest[1] = newTree
		forest = forest[1:]
	}
	return newTree
}

func idTree(t *tree, i int) {
	i++
	if t.l != nil {
		idTree(t.l, i)
	}
	t.id = i
	if t.r != nil {
		idTree(t.r, i)
	}
}

func treeToString(root *tree, s string) string {
	if root.l != nil {
		s = treeToString(root.l, s)
	}
	id := fmt.Sprintf("%08b", root.id)
	char := fmt.Sprintf("%08b", root.c)
	s = s + id + char
	if root.l != nil {
		left := fmt.Sprintf("%08b", root.l.id)
		s += left
	} else {
		left := fmt.Sprintf("%08b", 0)
		s += left
	}
	if root.r != nil {
		right := fmt.Sprintf("%08b", root.r.id)
		s += right
	} else {
		right := fmt.Sprintf("%08b", 0)
		s += right
	}
	if root.r != nil {
		s = treeToString(root.r, s)
	}
	return s
}

func storeLeaves(root *tree) []rune {
	s := make([]rune, 0)
	if root != nil {
		if root.l != nil {
			s = append(s, storeLeaves(root.l)...)
		}
		if root.l == nil && root.r == nil {
			s = append(s, root.c)
		}
		if root.r != nil {
			s = append(s, storeLeaves(root.r)...)
		}
	}
	return s
}

//traverse tree and store each value in order in some data structure
func traverse(root *tree) []rune {
	//Left, Node, Right
	s := make([]rune, 0)
	if root != nil {
		if root.l != nil {
			s = append(s, traverse(root.l)...)
		}
		s = append(s, root.c)
		if root.r != nil {
			s = append(s, traverse(root.r)...)
		}
	}
	return s
}

func paths(t *tree, path string) map[rune]string {
	if t == nil {
		return nil
	}
	result := merge(paths(t.l, path+"0"), paths(t.r, path+"1"))
	if t.c != 0 {
		result[t.c] = path
	}
	return result
}

func merge(maps ...map[rune]string) map[rune]string {
	result := make(map[rune]string)
	for _, m := range maps {
		for char, path := range m {
			result[char] = path
		}
	}
	return result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func makeTable(content string) map[rune]string {
	m := countChars(content)
	forest := makeForest(m)
	t := makeTree(forest)
	var path string
	table := paths(&t, path)
	return table
}

func stringToBits(s string, m map[rune]string) []byte {
	b := make([]byte, 1)
	eof := m['ൾ']
	var offset, i int
	for _, runeValue := range s {
		bitSequence := m[runeValue]
		for _, bit := range bitSequence {
			if offset == 8 {
				offset = 0
				var elem byte
				b = append(b, elem)
				i++
			}
			if bit == '0' {
				b[i] <<= 1
				offset++
			} else {
				b[i] <<= 1
				b[i] |= 1
				offset++
			}
		}
	}
	for _, bit := range eof {
		if offset == 8 {
			offset = 0
			i++
			var elem byte
			b = append(b, elem)
		}
		if bit == '0' {
			b[i] <<= 1
			offset++
		} else {
			b[i] <<= 1
			b[i] |= 1
			offset++
		}
	}
	b[i] <<= (8 - offset)
	return b
}

func treeStringToBits(s string) []byte {
	//fmt.Printf("%s\n", s)
	b := make([]byte, 1)
	var offset, i int
	for _, bit := range s {
		if offset == 8 {
			offset = 0
			var new byte
			b = append(b, new)
			i++
		}
		if bit == '0' {
			b[i] <<= 1
			offset++
		} else if bit == '1' {
			b[i] <<= 1
			b[i] |= 1
			offset++
		}
	}
	return b
}

func compress(fileName string) {
	f, err := os.Open(fileName)
	check(err)
	defer f.Close()
	b, err := ioutil.ReadFile(fileName)
	check(err)
	data := string(b) + "ൾ"
	m := countChars(data)
	forest := makeForest(m)
	t := makeTree(forest)
	var treeString string
	treeString := treeToString(t, treeString)
	compressedTree := treeStringToBits(treeString)
	var path string
	table := paths(&t, path)
	compressedFileName := fmt.Sprintf("%s.huff", fileName)
	cF, err := os.Create(compressedFileName)
	check(err)
	defer cF.Close()
	compressedData := append(compressedTree, stringToBits(data, table)...)
	n, err := cF.Write(compressedData)
	fmt.Printf("%v bytes written", n)
	check(err)
}

func bitsToTree(b []byte) []*tree {
	treeFields := make([]treeInfo, 0)
	m := make(map[int]*tree)
	for i := 0; i < len(b); i += 4 {
		var ti treeInfo
		var t tree
		ti.id = int(b[i])
		t.id = int(b[i])
		ti.c = int32(b[i+1])
		t.c = int32(b[i+1])
		ti.l = int(b[i+2])
		ti.r = int(b[i+3])
		treeFields = append(treeFields, ti)
		m[t.id] = &t
	}
	fmt.Printf("%v\n", m)
	fmt.Printf("%v\n", treeFields)
	newTree := make([]*tree, 0)
	for i := 0; i < len(treeFields); i++ {
		t := m[treeFields[i].id]
		if tl, ok := m[treeFields[i].l]; ok {
			t.l = tl
		}
		if tr, ok := m[treeFields[i].r]; ok {
			t.r = tr
		}
		newTree = append(newTree, t)
	}
	return newTree
}

func main() {
	/*content := "streets are stone stars are not"

	charCount := countChars(content)
	//fmt.Println(charCount)
	forest := makeForest(charCount)
	//fmt.Println(forest)
	t := makeTree(forest)
	var path string
	table := paths(&t, path)
	for c, p := range table {
		fmt.Printf("%c: %s \n", c, p)
	}
	*/
	fmt.Println("Enter the name of the file to compress:")
	var fileName string
	fmt.Scanln(&fileName)
	compress(fileName)
}
