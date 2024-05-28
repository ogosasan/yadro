package comics

import (
	"github.com/kljensen/snowball"
	"log"
	"strings"
	"unicode"
)

func Normalization(sentence string) []string {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && string(c) != "'"
	}
	words := strings.FieldsFunc(sentence, f)
	seen := make(map[string]interface{})
	res := []string{}

	for _, word := range words {
		stemmed, err := snowball.Stem(word, "english", false)
		if !IsStopWord(stemmed) {
			if err == nil && seen[stemmed] == nil {
				res = append(res, stemmed)
				seen[stemmed] = true
			} else if err != nil {
				log.Println(err)
				continue
			}
		}
	}
	return res
}
