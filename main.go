package main

import (
	"flag"
	"fmt"
	"github.com/kljensen/snowball"
	"strings"
)

func main() {
	var sentence string
	flag.StringVar(&sentence, "s", "default", "a string var")
	flag.Parse()
	words := strings.Fields(sentence)
	seen := make(map[string]bool)
	res := []string{}

	for _, word := range words {
		stemmed, err := snowball.Stem(word, "english", false)
		if !IsStopWord(stemmed) {
			if err == nil && !seen[stemmed] {
				res = append(res, stemmed)
				seen[stemmed] = true

			}
		}
	}
	fmt.Println(strings.Join(res, " "))
}
