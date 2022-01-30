package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type tree struct {
	id int
	c  rune
	w  int
	l  *tree
	r  *tree
}

type treeInfo struct {
	id int
	c  rune
	l  int
	r  int
}

func makeTree(content string) *tree {
	//countChars
	counts := make(map[rune]int)
	forest := make([]tree, len(content))
	var j int
	for _, c := range content {
		counts[c]++
		var charInForest bool
		for _, t := range forest {
			if c == t.c {
				charInForest = true
				break
			}
		}
		if charInForest == false {
			forest[j].c = c
			j++
		}
	}
	forest = forest[:len(counts)]
	for i := 0; i < len(forest); i++ {
		forest[i].w = counts[forest[i].c]
		forest[i].id = i + 1
	}
	//makeTree
	length := len(forest)
	var newTree tree
	for i := 0; len(forest) > 1; i++ {
		sort.Slice(forest, func(i, j int) bool { return forest[i].w < forest[j].w })
		newTree = combineTrees(forest[0], forest[1])
		newTree.id = length + i + 1
		forest[1] = newTree
		forest = forest[1:]
	}
	return &newTree
}

func combineTrees(t1 tree, t2 tree) tree {
	var t3 tree //= &tree{}
	t3.w = t1.w + t2.w
	t3.l = &t1
	t3.r = &t2
	return t3
}

func compressTree(root *tree, s string) string {
	if root.l != nil {
		s = compressTree(root.l, s)
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
		s = compressTree(root.r, s)
	}
	return s
}

func compressedTreeToBits(s string) []byte {
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
	treeEnd := []byte{1, 1, 1, 1, 0, 0, 0, 0}
	b = append(b, treeEnd...)
	return b
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

func stringToBits(s string, m map[rune]string) []byte {
	b := make([]byte, 1)
	//eof := m['ൾ']
	var offset, i int
	for _, runeValue := range s {
		bitSequence := m[runeValue]
		for _, bit := range bitSequence {
			if offset == 8 {
				offset = 0
				var new byte
				b = append(b, new)
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

	b[i] <<= (8 - offset)
	return b
}

func compress(fileName string) {
	f, err := os.Open(fileName)
	check(err)
	defer f.Close()
	b, err := ioutil.ReadFile(fileName)
	check(err)
	data := string(b) + "Þ"
	t := makeTree(data)
	var serialTree string
	serialTree = compressTree(t, serialTree)
	compressedTree := compressedTreeToBits(serialTree)
	fileName = strings.TrimSuffix(fileName, ".unhuff")
	compressedFileName := fmt.Sprintf("%s.huff", fileName)
	cF, err := os.Create(compressedFileName)
	check(err)
	defer cF.Close()
	var path string
	table := paths(t, path)
	compressedData := append(compressedTree, stringToBits(data, table)...)
	n, err := cF.Write(compressedData)
	fmt.Printf("%v bytes written", n)
	check(err)
}

func uncompressTree(b []byte) []*tree {
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

func findRoot(newTree []*tree) *tree {
	root := newTree[0]
	for _, t := range newTree {
		if t.l == root || t.r == root {
			root = t
		}
	}
	return root
}

func unhuff(data []byte, root *tree) string {
	var unhuffedData string
	trueRoot := root
	for _, b := range data {
		comp := byte(128)
		i := 0
		for i < 8 {
			if root.l != nil || root.r != nil {
				if b&comp == comp {
					root = root.r
				} else {
					root = root.l
				}
				i++
				comp >>= 1
			} else if root.c == 'Þ' {
				break
			} else {
				unhuffedData += string(root.c)
				root = trueRoot
			}
		}
	}
	return unhuffedData
}

func traverse(root *tree) []tree {
	//Left, Node, Right
	s := make([]tree, 0)
	if root != nil {
		if root.l != nil {
			s = append(s, traverse(root.l)...)
		}
		s = append(s, *root)
		if root.r != nil {
			s = append(s, traverse(root.r)...)
		}
	}
	return s
}

func findTreeEnd(b []byte) int {
	var oneByteCount int
	var oneSequenceDone bool
	var zeroByteCount int
	var treeEnd int
	for i, elem := range b {
		if elem == 1 {
			oneByteCount++
		} else if elem != 0 {
			oneByteCount = 0
		}
		if oneByteCount == 4 {
			oneSequenceDone = true
		} else {
			oneSequenceDone = false
		}
		if oneSequenceDone {
			if elem == 0 {
				zeroByteCount++
			} else {
				zeroByteCount = 0
			}
			if zeroByteCount == 4 {
				treeEnd = i - 7
			}
		}
	}
	return treeEnd
}

func decompress(fileName string) {
	if strings.HasSuffix(fileName, ".huff") {
		f, err := os.Open(fileName)
		check(err)
		defer f.Close()
		b, err := ioutil.ReadFile(fileName)
		check(err)
		treeEnd := findTreeEnd(b)
		fmt.Printf("tree end index: %v\n", treeEnd)
		treeData := b[:treeEnd]
		length := len(treeData)
		fmt.Printf("length of treeData: %v\n", length)
		fmt.Printf("sample of tree data: %v\n", treeData[length-5:])
		newTree := uncompressTree(treeData)
		fmt.Printf("nodes in tree: %v\n", len(newTree))
		root := findRoot(newTree)
		fmt.Printf("%p: %+v\n", root, root)
		data := b[treeEnd+8:]
		unhuffedData := unhuff(data, root)
		uncompressedFileName := strings.Replace(fileName, ".huff", ".unhuff", -1)
		uCF, err := os.Create(uncompressedFileName)
		check(err)
		defer uCF.Close()
		n, err := uCF.WriteString(unhuffedData)
		fmt.Printf("%v bytes written\n", n)
		check(err)
	} else {
		fmt.Printf("File is not in the correct format (.huff)\n")
	}
}

func main() {
	compress("testFile2")
	decompress("testFile2.huff")
}
