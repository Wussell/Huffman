package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type tree struct {
	char   rune
	weight int
	l      *tree
	r      *tree
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
		forest[i].char = c
		forest[i].weight = v
		fmt.Println(forest[i])
		i++
	}
	return forest
}

func combineTrees(t1 tree, t2 tree) tree {
	var t3 tree
	t3.weight = t1.weight + t2.weight
	t3.l = &t1
	t3.r = &t2
	return t3
}

func makeTree(forest []tree) tree {
	var newTree tree
	for len(forest) > 1 {
		sort.Slice(forest, func(i, j int) bool { return forest[i].weight < forest[j].weight })
		newTree = combineTrees(forest[0], forest[1])
		forest[1] = newTree
		forest = forest[1:]
	}
	return newTree
}

func storeLeaves(root *tree) []rune {
	s := make([]rune, 0)
	if root != nil {
		if root.l != nil {
			s = append(s, storeLeaves(root.l)...)
		}
		if root.l == nil && root.r == nil {
			s = append(s, root.char)
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
		s = append(s, root.char)
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
	if t.char != 0 {
		result[t.char] = path
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
	eof := "1111"
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
		} else if bit == '1' {
			b[i] <<= 1
			b[i] |= 1
			offset++
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
	data := string(b)
	table := makeTable(data)
	compressedData := stringToBits(data, table)
	compressedFileName := fmt.Sprintf("%s.huff", fileName)
	cF, err := os.Create(compressedFileName)
	check(err)
	defer cF.Close()
	n, err := cf.Write(compressedData)
	fmt.Printf("%v bytes written", n)
	check(err)
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
	compress("testFile")
}
