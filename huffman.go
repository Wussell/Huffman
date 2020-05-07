package main

import (
	"fmt"
	"sort"
)

type tree struct {
	char   rune
	weight int
	left   *tree
	right  *tree
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
	t3.left = &t1
	t3.right = &t2
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

/*
//combine the two trees of lowest weight in a forest until there is only one tree remaining
func makeTree(forest []tree) tree {
	var lowestWeight1, lw1I, lw2I, treeCount int = forest[0].weight, 0, 1, len(forest)
	var t3 tree
	for ; treeCount > 1; treeCount = len(forest) {
		//find index of two trees with lowest weights
		for i, t := range forest {
			if t.weight != 0 && t.weight < lowestWeight1 {
				lw2I = lw1I
				lowestWeight1 = t.weight
				lw1I = i
			}
		}
		t3 = combineTrees(forest[lw1I], forest[lw2I])
		//insert the new combined tree into the forest and remove the now extraneous tree
		forest[lw1I] = t3
		leftSlice := forest[:lw2I]
		if lw2I+1 < len(forest) {
			rightSlice := forest[lw2I+1:]
			for i := 0; i < len(rightSlice); i++ {
				leftSlice = append(leftSlice, rightSlice[i])
			}
		}
		forest = leftSlice
	}
	return t3
}

func traverse(t tree) []rune {
	current := &t
	if current.left != nil {

	}
}


func makeTable(t tree) map[rune]byte {

}
*/

func storeLeaves(root *tree) []rune {
	s := make([]rune, 0)
	if root != nil {
		if root.left != nil {
			s = append(s, storeLeaves(root.left)...)
		}
		if root.left == nil && root.right == nil {
			s = append(s, root.char)
		}
		if root.right != nil {
			s = append(s, storeLeaves(root.right)...)
		}
	}
	return s
}

func main() {
	content := "bookkeeping"
	charCount := countChars(content)
	fmt.Println(charCount)
	forest := makeForest(charCount)
	fmt.Println(forest)
	t := makeTree(forest)
	s := storeLeaves(&t)
	fmt.Println(s)
}
