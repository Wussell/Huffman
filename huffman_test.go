package main

import (
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
