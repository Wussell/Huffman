package huffman

type tree struct {
	char   rune
	weight int
	left   *tree
	right  *tree
}

func countChar(content string) map[rune]int {
	count := make(map[rune]int)
	for _, v := range content {
		count[v]++
	}
	return count
}

func createForest(charCount map[rune]int) []tree {
	var numChars int
	for range charCount {
		numChars++
	}
	forest := make([]tree, numChars)
	for _, tree := range forest {
		for c, v := range charCount {
			tree.char = c
			tree.weight = v
		}
	}
	return forest
}

func main() {

}
