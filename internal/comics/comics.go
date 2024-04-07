package comics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type GistRequest struct {
	Tscript string `json:"transcript"`
	Img     string `json:"img"`
}

type Write struct {
	Tscript []string `json:"keywords"`
	Img     string   `json:"url"`
}

type InfoResponse struct {
	Num int `json:"num"`
}

func GetNumComics(url string) (int, error) {
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var info InfoResponse
	err = json.NewDecoder(response.Body).Decode(&info)
	if err != nil {
		return 0, err
	}

	return info.Num, nil
}

func GoToSite(numComics int, baseURL string) map[int]Write {
	comicsMap := make(map[int]Write)

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(numComics)
	for i := 1; i <= numComics; i++ {
		go func(i int) {
			defer wg.Done()

			url := fmt.Sprintf(baseURL, i)

			response, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			if response.StatusCode != http.StatusOK {
				fmt.Printf("%s not found", url)
				fmt.Println(i)
				return
			}

			defer response.Body.Close()

			var xkcd GistRequest
			err = json.NewDecoder(response.Body).Decode(&xkcd)
			if err != nil {
				fmt.Println(err)
				return
			}

			xkcd.Tscript = normalization(xkcd.Tscript)
			print := Write{Tscript: strings.Fields(xkcd.Tscript), Img: xkcd.Img}

			mu.Lock()
			comicsMap[i] = print
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return comicsMap
}

func WriteFile(screen bool, num int, file string, comicsMap map[int]Write) {
	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	count := 0
	if screen {
		if num == -1 {
			bytes, _ := json.MarshalIndent(comicsMap, "", "\t")
			fmt.Println(string(bytes))
		} else {
			for _, gist := range comicsMap {
				if count >= num {
					break
				}
				bytes, _ := json.MarshalIndent(gist, "", "\t")
				fmt.Printf("%d:\n", count+1)
				fmt.Println(string(bytes))
				count++
			}
		}
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(comicsMap)
	if err != nil {
		fmt.Println(err)
		return
	}
}
