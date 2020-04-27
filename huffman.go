package main

import "fmt"

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
	var lowestWeight1, lowestWeight2, lw1I, lw2I, treeCount int = forest[0].weight, forest[1].weight, 0, 1, len(forest)
	var t3 tree
	for treeCount > 1 {
		for i, t := range forest {
			if t.weight != 0 && t.weight < lowestWeight1 {
				lowestWeight2 = lowestWeight1
				lw2I = lw1I
				lowestWeight1 = t.weight
				lw1I = i
			}
		}
		t3 = combineTrees(forest[lw1I], forest[lw2I])
		forest[lw1I] = t3
		leftSlice := forest[:lw2I]
		rightSlice := forest[lw2I+1:]
		forest = append(leftSlice, rightSlice)
		treeCount = len(forest)
	}
	return t3
}

func main() {
	content := "go go gophers"
	charCount := countChars(content)
	fmt.Println(charCount)
	forest := makeForest(charCount)
	fmt.Println(forest)
}
