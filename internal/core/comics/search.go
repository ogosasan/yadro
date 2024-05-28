package comics

import (
	"fmt"
	"sort"
)

type kv struct {
	Key   int
	Value int
}

func IndexSearch(indexMap map[string][]int, comicsMap map[int]Write, line string) []string {
	searchMap := make(map[int]int)
	words := Normalization(line)
	for i := 0; i < len(words); i++ {
		if indexMap[words[i]] != nil {
			for j := 0; j < len(indexMap[words[i]]); j++ {
				searchMap[indexMap[words[i]][j]]++ //делаю +1 по ключу j-го комикса, в котором есть слово i
			}
		}
	}
	//var ans []string
	ans := PrintUrl(searchMap, comicsMap)
	return ans
}

func PrintUrl(unsortedMap map[int]int, comicsMap map[int]Write) []string {
	var ss []kv
	var ans []string
	for k, v := range unsortedMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	count := 0
	for _, kv := range ss {
		ans = append(ans, comicsMap[kv.Key].Img)
		fmt.Println(comicsMap[kv.Key].Img)
		//fmt.Println(kv.Key)
		count++
		if count >= 10 {
			break
		}
	}
	return ans
}
