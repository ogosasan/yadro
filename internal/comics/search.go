package comics

import (
	"fmt"
	"sort"
)

type kv struct {
	Key   int
	Value int
}

func IndexSearch(indexMap map[string][]int, comicsMap map[int]Write, line string) {
	searchMap := make(map[int]int)
	words := normalization(line)
	for i := 0; i < len(words); i++ {
		if indexMap[words[i]] != nil {
			for j := 0; j < len(indexMap[words[i]]); j++ {
				searchMap[indexMap[words[i]][j]]++ //делаю +1 по ключу j-го комикса, в котором есть слово i
			}
		}
	}
	PrintUrl(searchMap, comicsMap)
}

func Search(comicsMap map[int]Write, str string) {
	wordsmap := make(map[int]int)
	words := normalization(str)
	for key, value := range comicsMap {
		for i := 0; i < len(value.Tscript); i++ {
			for j := 0; j < len(words); j++ {
				if value.Tscript[i] == words[j] {
					wordsmap[key]++
				}
			}
			{
			}
		}
	}
	PrintUrl(wordsmap, comicsMap)
}

func PrintUrl(ansortMap map[int]int, comicsMap map[int]Write) {
	var ss []kv
	for k, v := range ansortMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	count := 0
	for _, kv := range ss {
		fmt.Println(comicsMap[kv.Key].Img)
		count++
		if count >= 10 {
			break
		}
	}
}
