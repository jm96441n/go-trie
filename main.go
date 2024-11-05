package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type Node struct {
	Children [26]*Node
	EndNode  bool
}

func main() {
	fileContent, err := os.ReadFile("unique-words.md")
	if err != nil {
		log.Fatal(err)
	}

	words := strings.Split(string(fileContent), "\n")
	trie := buildTrie(words)

	if err := runServer(trie); err != nil {
		log.Fatal(err)
	}
}

func runServer(trie *Node) error {
	router := http.NewServeMux()
	router.HandleFunc("/search", trieHandler(trie))
	log.Println("Listening on port 8080")
	return http.ListenAndServe(":8080", router)
}

func buildTrie(words []string) *Node {
	root := &Node{Children: [26]*Node{}}

	for _, word := range words {
		insertWord(root, strings.ToLower(word), 0)
	}

	return root
}

func insertWord(node *Node, word string, idx int) {
	if idx == len(word) {
		node.EndNode = true
		return
	}

	char := word[idx]
	charIdx := char - 'a'
	if charIdx >= 26 {
		log.Fatalf("Invalid character: %c\nWord: %q", char, word)
	}
	child := node.Children[char-'a']
	if child == nil {
		child = &Node{Children: [26]*Node{}}
		node.Children[char-'a'] = child
	}
	insertWord(child, word, idx+1)
}

func searchWord(node *Node, word string, idx int) bool {
	if idx == len(word) {
		return node.EndNode
	}

	char := word[idx]
	child := node.Children[char-'a']
	if child == nil {
		return false
	}

	return searchWord(child, word, idx+1)
}

func trieHandler(trie *Node) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		word := r.URL.Query().Get("word")
		if word == "" {
			http.Error(w, "missing word query parameter", http.StatusBadRequest)
			return
		}
		if searchWord(trie, strings.ToLower(word), 0) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Word found"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Word not found"))
		}
	}
}
