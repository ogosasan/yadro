package comics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

func GetNumComics(baseUrl string) int {
	i := 1000
	j := 500
	for {
		url := fmt.Sprintf(baseUrl, i)
		response, _ := http.Get(url)
		if response.StatusCode == http.StatusOK {
			i += j
			continue
		} else {
			i -= j
			j = j / 2
		}
		if j == 0 {
			break
		}
	}
	return i
}

func GoToSite(numComics int, baseURL string, done chan os.Signal, fileExist bool, workers int) (map[int]Write, map[string][]string, int) {
	existComics := make(map[int]Write)
	if fileExist {
		data, err := ioutil.ReadFile("database.json")
		if err != nil {
			fmt.Print(err)
		}
		err = json.Unmarshal(data, &existComics)
		if err != nil {
			fmt.Println("error:", err)
		}
	}

	comicsMap := make(map[int]Write)
	indexMap := make(map[string][]string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	jobs := make(chan int, numComics)
	results := make(chan Write, numComics)
	wg.Add(workers)

	worker := func() {
		defer wg.Done()
		for i := range jobs {
			url := fmt.Sprintf(baseURL, i)
			response, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if response.StatusCode != http.StatusOK {
				fmt.Printf("%s not found", url)
				continue
			}
			var xkcd GistRequest
			err = json.NewDecoder(response.Body).Decode(&xkcd)
			if err != nil {
				fmt.Println(err)
				continue
			}
			printInFile := Write{Tscript: normalization(xkcd.Tscript), Img: xkcd.Img}
			results <- printInFile
			response.Body.Close()
		}
	}

	for i := 0; i < workers; i++ {
		go worker()
	}

	for i := 1; i <= numComics; i++ {
		if existComics[i].Img == "" && len(existComics[i].Tscript) == 0 {
			jobs <- i
		} else {
			comicsMap[i] = existComics[i]
		}
	}
	close(jobs)

	go func() {
		for result := range results {
			mu.Lock()
			comicsMap[len(comicsMap)+1] = result
			mu.Unlock()
		}
	}()

	go func() {
		wg.Wait()
		close(results)
		close(done)
	}()
	wg.Wait()
	for key, value := range comicsMap {
		for j := 0; j < len(value.Tscript); j++ {
			indexMap[value.Tscript[j]] = append(indexMap[value.Tscript[j]], strconv.Itoa(key))
		}
	}
	return comicsMap, indexMap, len(existComics)
}

func WriteFile(file string, comicsMap map[int]Write, indexMap map[string][]int) {
	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(comicsMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	index, err := os.Create("index.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer index.Close()
	encode := json.NewEncoder(index)
	encode.SetIndent("", "\t")
	err = encode.Encode(indexMap)
	if err != nil {
		fmt.Println(err)
		return
	}
}
