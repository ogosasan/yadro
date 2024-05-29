package comics

import (
	"reflect"
	"testing"
)

func TestIndexSearch(t *testing.T) {
	indexMap := map[string][]int{
		"boy":      {1},
		"sit":      {1},
		"barrel":   {1},
		"float":    {1},
		"ocean":    {1},
		"wonder":   {1},
		"next":     {1},
		"drift":    {1},
		"distance": {1},
		"nothing":  {1},
		"else":     {1},
		"seen":     {1},
		"alt":      {1},
	}

	comicsMap := map[int]Write{
		1: {
			Tscript: []string{"boy", "sit", "barrel", "float", "ocean", "wonder", "next", "drift", "distance", "nothing", "else", "seen", "alt"},
			Img:     "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg",
		},
	}

	line := "boy sit"
	expected := []string{"https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg"}

	result := IndexSearch(indexMap, comicsMap, line)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestPrintUrl(t *testing.T) {
	ansortMap := map[int]int{
		1: 2,
		2: 1,
	}

	comicsMap := map[int]Write{
		1: {
			Tscript: []string{"boy", "sit", "barrel", "float", "ocean", "wonder", "next", "drift", "distance", "nothing", "else", "seen", "alt"},
			Img:     "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg",
		},
		2: {
			Tscript: []string{"cat", "dog"},
			Img:     "https://imgs.xkcd.com/comics/(2).jpg",
		},
	}

	expected := []string{"https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg", "https://imgs.xkcd.com/comics/(2).jpg"}

	result := PrintUrl(ansortMap, comicsMap)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
