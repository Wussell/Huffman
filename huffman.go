package main

import (
	"fmt"
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

func main() {
	content := "streets are stone stars are not"
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
}
