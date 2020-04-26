package huffman

func charCount(content string) map[string]int {
	count := make(map[string]int)
	for _, v := range content {
		count[v]++
	}
}

func main() {

}
