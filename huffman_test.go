package main

import (
	"fmt"
	"testing"
)

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

resultSample := make([]byte, 10)

func BenchmarkCompress(b *testing.B) {
	data := []byte("We the People of the United States, in Order to form a more perfect Union, establish Justice, insure domestic Tranquility, provide for the common defense, promote the general Welfare, and secure the Blessings of Liberty to ourselves and our Posterity, do ordain and establish this Constitution for the United States of America")
	d := make([]byte, len(data))
	for n := 0; n < b.N; n++ {
		d = compress(data)
	}
	fmt.Println(d[0:10])
	resultSample = d[0:10]
}

var store string 

func BenchmarkDecompress(b *testing.B) {
	data := []byte{21, 109, 0, 0, 51, 0, 21, 25, 25, 117, 0, 0, 57, 0, 51, 52, 18, 44, 0, 0, 52, 0, 18, 47, 31, 66, 0, 0, 39, 0, 31, 19, 19, 79, 0, 0, 41, 0, 39, 1, 1, 87, 0, 0, 47, 0, 41, 44, 29, 118, 0, 0, 44, 0, 29, 6, 6, 80, 0, 0, 63, 0, 57, 58, 26, 84, 0, 0, 36, 0, 26, 34, 34, 65, 0, 0, 43, 0, 36, 30, 30, 103, 0, 0, 48, 0, 43, 42, 15, 83, 0, 0, 42, 0, 15, 37, 27, 113, 0, 0, 37, 0, 27, 33, 33, 67, 0, 0, 53, 0, 48, 9, 9, 108, 0, 0, 58, 0, 53, 12, 12, 110, 0, 0, 67, 0, 63, 64, 5, 104, 0, 0, 54, 0, 5, 10, 10, 102, 0, 0, 59, 0, 54, 17, 17, 115, 0, 0, 64, 0, 59, 2, 2, 101, 0, 0, 69, 0, 67, 68, 13, 105, 0, 0, 60, 0, 13, 20, 20, 114, 0, 0, 65, 0, 60, 61, 8, 112, 0, 0, 49, 0, 8, 46, 11, 85, 0, 0, 46, 0, 11, 23, 23, 98, 0, 0, 55, 0, 49, 14, 14, 100, 0, 0, 61, 0, 55, 7, 7, 111, 0, 0, 68, 0, 65, 66, 3, 32, 0, 0, 66, 0, 3, 62, 22, 99, 0, 0, 50, 0, 22, 45, 35, 222, 0, 0, 40, 0, 35, 38, 24, 74, 0, 0, 38, 0, 24, 32, 32, 76, 0, 0, 45, 0, 40, 28, 28, 121, 0, 0, 56, 0, 50, 16, 16, 97, 0, 0, 62, 0, 56, 4, 4, 116, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 26, 247, 161, 225, 238, 232, 43, 214, 157, 232, 122, 145, 199, 186, 226, 95, 223, 106, 45, 7, 134, 102, 174, 119, 222, 77, 200, 55, 112, 46, 94, 161, 202, 95, 31, 169, 28, 89, 139, 53, 254, 211, 44, 42, 55, 40, 43, 241, 195, 22, 131, 80, 203, 213, 176, 53, 248, 227, 16, 79, 76, 152, 48, 88, 254, 98, 212, 77, 142, 138, 188, 155, 157, 232, 123, 139, 0, 44, 245, 90, 89, 171, 22, 162, 108, 23, 239, 122, 30, 34, 205, 207, 75, 134, 178, 167, 178, 197, 186, 117, 202, 248, 12, 189, 232, 120, 96, 173, 86, 12, 138, 235, 78, 229, 197, 55, 63, 207, 125, 235, 12, 171, 40, 230, 187, 167, 93, 97, 156, 62, 215, 220, 199, 243, 22, 173, 235, 154, 246, 15, 116, 235, 154, 255, 105, 150, 21, 27, 209, 11, 137, 236, 215, 227, 195, 241, 103, 147, 115, 189, 15, 82, 56, 247, 92, 75, 251, 237, 117, 167, 16, 129, 204, 113, 222, 64}
	var s string
	for n := 0; n < b.N; n++ {
		s = decompress(data)
	}
	fmt.Println(s)
	store = s
}
